package wikipedia

import (
	"encoding/json"
	"fmt"
)

// Links is a from-to mapping of directional links using page titles.
type Links map[string][]string

func (pl Links) add(from, to string) {
	if _, ok := pl[from]; !ok {
		pl[from] = []string{}
	}
	pl[from] = append(pl[from], to)
}

func LinksFrom(titles []string) chan Links {
	return allLinks("pl", "links", titles)
}

func LinksHere(titles []string) chan Links {
	return allLinks("lh", "linkshere", titles)
}

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
					// TODO exponential backoff?
					panic(err)
				}

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

// extractContinue gets the continue string from the structure:
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
	pl := Links{}

	query := data["query"].(map[string]interface{})
	pages := query["pages"].(map[string]interface{})
	for _, page := range pages {
		pageMap := page.(map[string]interface{})
		fromTitle := pageMap["title"].(string)

		linksSlice, ok := pageMap[subkey].([]interface{})
		if ok {
			for _, link := range linksSlice {
				linkMap := link.(map[string]interface{})
				pl.add(fromTitle, linkMap["title"].(string))
			}
		}
	}

	return pl
}
