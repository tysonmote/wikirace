package wikipedia

import (
	"encoding/json"
	"reflect"
	"testing"
)

const (
	partialLinksJSON = `{
		"continue": {
				"plcontinue": "39027|0|AAU_Junior_Olympic_Games",
				"continue": "||"
		},
		"query": {
				"pages": {
						"7365423": {
								"pageid": 7365423,
								"ns": 0,
								"title": "Tryall Golf Club"
						},
						"39027": {
								"pageid": 39027,
								"ns": 0,
								"title": "Mike Tyson",
								"links": [
										{
												"ns": 0,
												"title": "1984 Summer Olympics"
										},
										{
												"ns": 0,
												"title": "20/20 (US television show)"
										},
										{
												"ns": 0,
												"title": "2009 Golden Globe Awards"
										}
								]
						}
				}
		},
		"limits": {
				"links": 3
		}
	}`
)

func TestLinksResponse_UnmarshalJSON(t *testing.T) {
	resp := linksResponse{prefix: "pl", prop: "links"}
	err := json.Unmarshal([]byte(partialLinksJSON), &resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Continue != "39027|0|AAU_Junior_Olympic_Games" {
		t.Errorf("unexpected continute: %#v", resp.Continue)
	}

	expectLinks := Links{
		"Mike Tyson": []string{
			"1984 Summer Olympics",
			"20/20 (US television show)",
			"2009 Golden Globe Awards",
		},
	}

	if !reflect.DeepEqual(expectLinks, resp.Links) {
		t.Errorf("expected: %#v\ngot: %#v", expectLinks, resp.Links)
	}
}
