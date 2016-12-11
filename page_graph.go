package main

import (
	"fmt"

	"github.com/tysontate/wikirace/wikipedia"
)

// -- PageGraph

type PageGraph struct {
	root *Page
	all  map[string]*Page
}

func NewPageGraph(rootTitle string) PageGraph {
	root := NewPage(rootTitle)
	root.parent = NewPage("")
	return PageGraph{
		root: root,
		all:  map[string]*Page{rootTitle: root},
	}
}

func (g *PageGraph) UnvisitedLinks(title string) []string {
	page := g.all[title]
	links := page.Links()
	unvisitedLinks := []string{}

	for _, link := range links {
		linkedPage := g.getOrAdd(link)
		if linkedPage.parent == nil {
			linkedPage.parent = page
			unvisitedLinks = append(unvisitedLinks, link)
		}
	}
	return unvisitedLinks
}

func (g *PageGraph) Path(title string) []string {
	path := []string{}
	parent := g.all[title]
	for parent != nil && parent != g.root {
		path = append(path, parent.title)
		parent = parent.parent
	}
	path = append(path, g.root.title)

	pathLen := len(path)
	for i := 0; i < pathLen/2; i++ {
		swap := pathLen - i - 1
		path[i], path[swap] = path[swap], path[i]
	}

	return path
}

func (g *PageGraph) getOrAdd(title string) *Page {
	page := g.all[title]
	if page == nil {
		page = NewPage(title)
		g.all[title] = page
	}
	return page
}

// -- Page

type Page struct {
	title  string
	parent *Page // parent with shortest path to root of graph
	links  map[string]bool
}

func NewPage(title string) *Page {
	return &Page{
		title: title,
	}
}

func (p *Page) Links() []string {
	if p.links == nil {
		p.links = map[string]bool{}
		for links := range wikipedia.LinksFrom([]string{p.title}) {
			for title, linkTitles := range links {
				if title != p.title {
					panic(fmt.Errorf("expected links for %#v, got links for %#v", p.title, title))
				}
				for _, linkTitle := range linkTitles {
					p.links[linkTitle] = true
				}
			}
		}
	}

	links := []string{}
	for link := range p.links {
		links = append(links, link)
	}
	return links
}
