package api

import (
	"reflect"
	"strings"
	"testing"
)

func TestBatch(t *testing.T) {
	tests := []struct {
		given  []string
		size   int
		expect [][]string
	}{
		{[]string{}, 3, [][]string{}},
		{[]string{"a"}, 3, [][]string{{"a"}}},
		{[]string{"a", "b", "c"}, 3, [][]string{{"a", "b", "c"}}},
		{[]string{"a", "b", "c", "d"}, 3, [][]string{{"a", "b", "c"}, {"d"}}},
		{[]string{"a", "b", "c", "d", "e"}, 2, [][]string{{"a", "b"}, {"c", "d"}, {"e"}}},
	}

	for i, test := range tests {
		got := batch(test.given, test.size)
		if !reflect.DeepEqual(test.expect, got) {
			t.Errorf("tests[%d]: expected: %#v, got: %#v", i, test.expect, got)
		}
	}
}

func TestBuildQueryURL(t *testing.T) {
	url := buildQueryURL("xx", "titles", []string{"foo", "bar"}, "abc")

	params := []string{
		"prop=titles",
		"titles=foo%7Cbar",
		"xxcontinue=abc",
		"xxlimit=max",
		"xxnamespace=0",
	}
	for _, expected := range params {
		if !strings.Contains(url, expected) {
			t.Errorf("expected to find %#v in %#v", expected, url)
		}
	}
}
