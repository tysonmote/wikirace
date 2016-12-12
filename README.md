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

wikirace has no external dependencies. Just build with `go build`

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
