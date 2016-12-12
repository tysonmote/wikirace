package main

import (
	"log"
	"sync"

	"github.com/tysontate/wikirace/wikipedia"
)

// PageGraph represents a graph of Wikipedia pages that is built using
// a bidirectional breadth-first search (forwards from a starting page and
// backwards from an ending page).
//
// While searching, PageGraph runs two goroutines and, therefore, will at most
// have two simultaneous API requests running against Wikipedia at a time.
type PageGraph struct {
	// map of page titles to their parent page title
	forward map[string]string

	// queue of pages to search forwards from
	forwardQueue []string

	// map of page titles to their child page title
	backward map[string]string

	// queue of pages to search backwards from
	backwardQueue []string

	// For simplicity's sake, all maps share a lock. If lock contention is
	// a problem, we can split this in to two locks.
	sync.RWMutex
}

func NewPageGraph() PageGraph {
	return PageGraph{
		forward:       map[string]string{},
		forwardQueue:  []string{},
		backward:      map[string]string{},
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
		path = append(path, cursor)
		cursor = pg.forward[cursor]
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
		path = append(path, cursor)
		cursor = pg.backward[cursor]
	}

	return path
}

// Returns midpoint node, if full path is found
func (pg *PageGraph) SearchForward(from string) string {
	pg.forwardQueue = append(pg.forwardQueue, from)

	for len(pg.forwardQueue) != 0 {
		pages := pg.forwardQueue
		pg.forwardQueue = []string{}
		for links := range wikipedia.LinksFrom(pages) {
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
	pg.RLock()
	_, ok := pg.forward[to]
	pg.RUnlock()
	if !ok {
		log.Printf("FORWARD %#v -> %#v", from, to)
		// "to" page doesn't have a path to the source yet.
		pg.Lock()
		pg.forward[to] = from
		pg.forwardQueue = append(pg.forwardQueue, to)
		pg.Unlock()
	}

	// If we now have a path to the destination, we're done!
	pg.RLock()
	_, done = pg.backward[to]
	pg.RUnlock()

	return done
}

// Returns midpoint node, if full path is found
func (pg *PageGraph) SearchBackward(to string) string {
	pg.backwardQueue = append(pg.backwardQueue, to)

	for len(pg.backwardQueue) != 0 {
		pages := pg.backwardQueue
		pg.backwardQueue = []string{}
		for links := range wikipedia.LinksFrom(pages) {
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
	pg.RLock()
	_, ok := pg.backward[from]
	pg.RUnlock()
	if !ok {
		log.Printf("BACKWARD %#v -> %#v", from, to)
		// "from" page doesn't have a path to the destination yet.
		pg.Lock()
		pg.backward[from] = to
		pg.backwardQueue = append(pg.backwardQueue, from)
		pg.Unlock()
	}

	// If we now have a path to the source, we're done!
	pg.RLock()
	_, done = pg.forward[to]
	pg.RUnlock()
	return done
}
