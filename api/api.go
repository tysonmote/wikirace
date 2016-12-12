package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	apiEndpoint = "https://en.wikipedia.org/w/api.php"

	// Per https://meta.wikimedia.org/wiki/User-Agent_policy
	userAgent = "wikirace/0.1 (https://github.com/tysontate/wikirace; tyson@bufio.net)"

	// We only care about primary Wikipedia pages, not meta pages like "Talk",
	// etc. All namespaces are documented here:
	// https://en.wikipedia.org/wiki/Wikipedia:Namespace
	namespace = "0|14|100"
)

var (
	client = &http.Client{Timeout: 5 * time.Second}
)

// batch returns the given slice as batches with a maximum size.
func batch(s []string, max int) [][]string {
	batches := [][]string{}
	var start, end int
	for start < len(s) {
		end = start + max
		if end > len(s) {
			end = len(s)
		}
		batches = append(batches, s[start:end])
		start = end
	}
	return batches
}

// buildQueryURL builds a Wikipedia-style GET API request URL string for
// querying "links" or "linkshere" properties.
func buildQueryURL(prefix, prop string, titles []string, cont string) string {
	values := url.Values{}
	values.Add("format", "json")
	values.Add("action", "query")
	values.Add("titles", strings.Join(titles, "|"))
	values.Add("prop", prop)
	values.Add(fmt.Sprintf("%snamespace", prefix), namespace)
	values.Add(fmt.Sprintf("%slimit", prefix), "max")
	if len(cont) > 0 {
		values.Add(fmt.Sprintf("%scontinue", prefix), cont)
	}
	return fmt.Sprintf("%s?%s", apiEndpoint, values.Encode())
}

// get runs a Wikipedia GET API call and returns the full response body or an
// error if the request was not successful.
func get(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got status code: %s", resp.Status)
	}

	return ioutil.ReadAll(resp.Body)
}
