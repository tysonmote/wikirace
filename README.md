# wikirace

Given two Wikipedia pages, wikirace will try to find a short path from one page
to another using only links to other Wikipedia pages.

wikirace uses live data (i.e. no pre-calculated link graphs) but is still
extremely fast, even when using Wikipedia's list of [difficult-to-reach
pages][good_target]:

```
% time ./wikirace Altoids Doorbell
Altoids
Cinnamon
China
Door
Doorbell
./wikirace Altoids Doorbell  0.11s user 0.03s system 17% cpu 0.776 total

% time ./wikirace "Preparation H" Pokémon
Preparation H
The New York Times
Chicago Sun-Times
Pokémon
./wikirace "Preparation H" Pokémon  0.10s user 0.03s system 17% cpu 0.718 total
```

## Building

wikirace has no external dependencies. Just fetch and build with: `go get
github.com/tysontate/wikirace`

## Running

```
usage: ./wikirace [-debug] from_title to_title

  -debug
      print debugging log output to stderr
```

Example:

```
% wikirace "Mike Tyson" "Oceanography"
Mike Tyson
Alexander the Great
Aegean Sea
Oceanography
```

[good_target]: https://en.wikipedia.org/wiki/Wikipedia:Wikirace#Good_Target_Pages

## Limitations

* wikirace adheres to the [WikiMedia etiquette guide][etiquette] as faithfully
  as possible. To that end, it runs, at most, two simultaneous API requests to
  Wikipedia at a time.

* Wikipedia's API for fetching links from / to pages isn't as granular as the
  raw HTML, which can make it hard to exclude "boring" link paths. For example,
  many pages have an "[Authority control][auth_control]" block which has links
  to pages like "International Standard Book Number" which are linked to from
  other pages with "Authority control" sections. I've excluded as many of those
  as I could find.

[etiquette]: https://www.mediawiki.org/wiki/API:Etiquette
[auth_control]: https://en.wikipedia.org/wiki/Help:Authority_control

## Process

* **2 hours** - Researching possible strategies for building wikirace including
  tools, libraries, algorithms, and Wikipedia's requirements for programmatic
  access to their site. I considered screen scraping, but that's too laborious
  and error-prone given that Wikipedia offers a full API. I considered lots of
  different search algorithms, but a basic bidirectional depth-first search
  seemed to make the most sense given that Wikipedia's API allows you to query
  links both from and *to* a given page, making the backwards part of
  bidirectional search possible. Bidirectional search would also (later on in my
  process) allow me to batch requests together so that I could minimize the
  number of requests.

* **2 hours** - Writing basic Wikipedia API code. Wikipedia's API is somewhat
  unusual / bespoke and has lots of little quirks that I had to find through
  trial and error. I ended up rewriting my API code a couple times to avoid all
  the duplication of code that my initial versions had.

* **1 hour** - First pass at a basic unidirectional breadth-first search that
  made an API request for every graph node visit. I worked out lots of kinks
  with my Wikipedia API code and ended up with an extremely slow but working
  implementation of wikirace that I could play with.

* **2 hours** - Converted the one-request-per-page API code to a batch model
  that would allow me to fetch (for example) the adjacent pages for a list of 50
  pages at one time rather than issuing an API request for each and every page.
  This model adheres to Wikipedia's published API limits.

* **2 hours** - Replaced unidirectional depth-first search with a bidirectional
  depth-first search. This is orders of magnitude faster in practice -- I'm able
  to find short paths between unrelated pages in around a second from my local
  machine. Car -> Petunia in 0.8s, Mike Tyson -> Carp in 0.7s, Pencil -> Calcium
  in 1.2s, Google -> Wheat in 0.7s, etc.

* **1 hour** - Writing README, tidying up some odds and ends, adding some
  documentation throughout.
