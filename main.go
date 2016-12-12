package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var (
	// flags
	debug = flag.Bool("debug", false, "print debugging log output to stderr")

	// arguments
	fromTitle string
	toTitle   string
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [-debug] from_title to_title\n\n", os.Args[0])
	flag.PrintDefaults()
}

func init() {
	flag.Usage = usage
	flag.Parse()

	if !*debug {
		log.SetOutput(ioutil.Discard)
	}

	fromTitle = flag.Arg(0)
	toTitle = flag.Arg(1)

	if len(fromTitle) == 0 || len(toTitle) == 0 {
		usage()
		os.Exit(1)
	}
}

func main() {
	graph := NewPageGraph()
	for _, page := range graph.Search(fromTitle, toTitle) {
		fmt.Println(page)
	}
}
