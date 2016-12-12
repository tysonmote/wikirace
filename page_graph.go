package main

import (
	"log"
	"sync"

	"github.com/tysontate/wikirace/api"
)

type safeStringMap struct {
	strings map[string]string
	sync.RWMutex
}

func newSafeStringMap() safeStringMap {
	return safeStringMap{map[string]string{}, sync.RWMutex{}}
}

func (m *safeStringMap) Get(key string) (value string, exists bool) {
	m.RLock()
	defer m.RUnlock()
	value, exists = m.strings[key]
	return
}

func (m *safeStringMap) Set(key, value string) {
	m.Lock()
	defer m.Unlock()
	m.strings[key] = value
}

// PageGraph represents a graph of Wikipedia pages that is built using
// a bidirectional breadth-first search (forwards from a starting page and
// backwards from an ending page).
//
// While searching, PageGraph runs two goroutines and, therefore, will at most
// have two simultaneous API requests running against Wikipedia at a time.
type PageGraph struct {
	// map of page titles to their parent page title
	forward safeStringMap

	// queue of pages to search forwards from
	forwardQueue []string

	// map of page titles to their child page title
	backward safeStringMap

	// queue of pages to search backwards from
	backwardQueue []string
}

func NewPageGraph() PageGraph {
	return PageGraph{
		forward:       newSafeStringMap(),
		forwardQueue:  []string{},
		backward:      newSafeStringMap(),
		backwardQueue: []string{},
	}
}

// Search takes starting and ending page titles and returns a short path of
// links from the starting page to the ending page.
func (pg *PageGraph) Search(from, to string) []string {
	midpoint := make(chan string)

	go func() { midpoint <- pg.SearchForward(from) }()
	go func() { midpoint <- pg.SearchBackward(to) }()

	return pg.path(<-midpoint)
}

func (pg *PageGraph) path(midpoint string) []string {
	path := []string{}

	// Build path from start to midpoint
	cursor := midpoint
	for len(cursor) > 0 {
		log.Printf("FOUND PATH FORWARD: %#v", cursor)
		path = append(path, cursor)
		cursor, _ = pg.forward.Get(cursor)
	}
	for i := 0; i < len(path)/2; i++ {
		swap := len(path) - i - 1
		path[i], path[swap] = path[swap], path[i]
	}

	// Pop off midpoint because following loop adds it back in
	path = path[0 : len(path)-1]

	// Add path from midpoint to end
	cursor = midpoint
	for len(cursor) > 0 {
		log.Printf("FOUND PATH BACKWARDS: %#v", cursor)
		path = append(path, cursor)
		cursor, _ = pg.backward.Get(cursor)
	}

	return path
}

// Returns midpoint node, if full path is found
func (pg *PageGraph) SearchForward(from string) string {
	pg.forward.Set(from, "")
	pg.forwardQueue = append(pg.forwardQueue, from)

	for len(pg.forwardQueue) != 0 {
		pages := pg.forwardQueue
		pg.forwardQueue = []string{}
		log.Printf("SEARCHING FORWARD: %#v", pages)
		for links := range api.LinksFrom(pages) {
			for from, tos := range links {
				for _, to := range tos {
					if pg.checkForward(from, to) {
						return to
					}
				}
			}
		}
	}

	log.Println("forward queue is empty, returning")
	return ""
}

func (pg *PageGraph) checkForward(from, to string) (done bool) {
	_, exists := pg.forward.Get(to)
	if !exists {
		log.Printf("FORWARD %#v -> %#v", from, to)
		// "to" page doesn't have a path to the source yet.
		pg.forward.Set(to, from)
		pg.forwardQueue = append(pg.forwardQueue, to)
	}

	// If we now have a path to the destination, we're done!
	_, done = pg.backward.Get(to)
	return done
}

// Returns midpoint node, if full path is found
func (pg *PageGraph) SearchBackward(to string) string {
	pg.backward.Set(to, "")
	pg.backwardQueue = append(pg.backwardQueue, to)

	for len(pg.backwardQueue) != 0 {
		pages := pg.backwardQueue
		pg.backwardQueue = []string{}
		log.Printf("SEARCHING BACKWARD: %#v", pages)
		for links := range api.LinksFrom(pages) {
			for to, froms := range links {
				for _, from := range froms {
					if pg.checkBackward(from, to) {
						return to
					}
				}
			}
		}
	}

	log.Println("backward queue is empty, returning")
	return ""
}

func (pg *PageGraph) checkBackward(from, to string) (done bool) {
	_, exists := pg.backward.Get(from)
	if !exists {
		log.Printf("BACKWARD %#v -> %#v", from, to)
		// "from" page doesn't have a path to the destination yet.
		pg.backward.Set(from, to)
		pg.backwardQueue = append(pg.backwardQueue, from)
	}

	// If we now have a path to the source, we're done!
	_, done = pg.forward.Get(to)
	return done
}
