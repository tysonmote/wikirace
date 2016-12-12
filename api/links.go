package api

import (
	"encoding/json"
	"fmt"
)

var (
	// These page titles lead to boring wikirace path results and are, thus,
	// excluded.
	boring = map[string]bool{
		// Authority control section
		"Bibliothèque nationale de France":       true,
		"Digital object identifier":              true,
		"Integrated Authority File":              true,
		"International Standard Book Number":     true,
		"International Standard Name Identifier": true,
		"Library of Congress Control Number":     true,
		"MusicBrainz":                            true,
		"National Diet Library":                  true,
		"Virtual International Authority File":   true,
	}
)

// Links is a mapping of directional page links using page titles.
type Links map[string][]string

func (pl Links) add(from, to string) {
	if boring[from] || boring[to] {
		return
	}

	// The API can return pages that link to themselves. We should ignore them.
	if from == to {
		return
	}

	if _, ok := pl[from]; !ok {
		pl[from] = []string{}
	}
	pl[from] = append(pl[from], to)
}

// LinksFrom takes one or more Wikipedia page titles and returns a channel that
// will receive one or more Links objects, each containing partial or full
// mappings of page to linked page. The channel will be closed after all
// results have been fetched.
func LinksFrom(titles []string) chan Links {
	return allLinks("pl", "links", titles)
}

// LinksFrom takes one or more Wikipedia page titles and returns a channel that
// will receive one or more Links objects, each containing partial or full
// mappings of linked page to source page. The channel will be closed after all
// results have been fetched.
func LinksHere(titles []string) chan Links {
	return allLinks("lh", "linkshere", titles)
}

// allLinks batches API requests to fetch the maximum number of results allowed
// by Wikipedia and then sends Links objects containing those responses from
// Wikipedia on the returned channel.
func allLinks(prefix, prop string, titles []string) chan Links {
	c := make(chan Links)

	go func(prefix, prop string, titles []string) {
		// Holds Wikipedia's "continue" string if we have more results to fetch.
		// Set after the first request.
		var cont string

		// Wikipedia can batch process up to 50 page titles at a time.
		for _, titlesBatch := range batch(titles, 50) {
			// Continue paginating through results as long as Wikipedia is telling us
			// to continue.
			for i := 0; i == 0 || len(cont) > 0; i++ {
				queryURL := buildQueryURL(prefix, prop, titlesBatch, cont)
				body, err := get(queryURL)
				if err != nil {
					// If Wikipedia returns an error, just panic instead of doing an
					// exponential back-off.
					panic(err)
				}

				// Parse the response.
				resp := linksResponse{prefix: prefix, prop: prop}
				err = json.Unmarshal(body, &resp)
				if err != nil {
					panic(err)
				}

				c <- resp.Links
				cont = resp.Continue
			}
		}
		close(c)
	}(prefix, prop, titles)

	return c
}

// -- api response format

// linksResponse encapsulates Wikipedia's query API response with either
// "links" or "linkshere" properties enumerated.
type linksResponse struct {
	prefix   string
	prop     string
	Continue string
	Links    Links
}

func (r *linksResponse) UnmarshalJSON(b []byte) error {
	data := map[string]interface{}{}
	json.Unmarshal(b, &data)

	r.Continue = extractContinue(data, fmt.Sprintf("%scontinue", r.prefix))
	r.Links = extractLinks(data, r.prop)

	return nil
}

// extractContinue takes as input a Wikipedia API query response and returns
// the "continue" string. If no continue string is set, an empty string is
// returned.
//
//   {
//     "continue": {
//       "{subkey}": "736|0|Action-angle_variables"
//     }
//   }
func extractContinue(data map[string]interface{}, subkey string) string {
	if cont, ok := data["continue"]; ok {
		if contValue, ok := cont.(map[string]interface{})[subkey]; ok {
			return contValue.(string)
		}
	}
	return ""
}

// extractLinks takes as input a Wikipedia API query response with either
// "links" or "linkshere" properties enumerated for a set of pages and returns
// a complete Links representation of that response.
//
//   {
//     ...
//     "query": {
//       "pages": {
//         "15580374": {
//           "title": "Albert Einstein",
//           "{subkey}":[
//             { "title": "2dF Galaxy Redshift Survey" },
//             ...
//           ]
//         },
//         ...
//       }
//     }
//   }
func extractLinks(data map[string]interface{}, subkey string) Links {
	links := Links{}

	query := data["query"].(map[string]interface{})
	pages := query["pages"].(map[string]interface{})
	for _, page := range pages {
		pageMap := page.(map[string]interface{})
		fromTitle := pageMap["title"].(string)

		linksSlice, ok := pageMap[subkey].([]interface{})
		if ok {
			for _, link := range linksSlice {
				linkMap := link.(map[string]interface{})
				links.add(fromTitle, linkMap["title"].(string))
			}
		}
	}

	return links
}
