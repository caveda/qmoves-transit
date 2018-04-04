package transit

import (
	"os"
	"testing"
)

var getSourcesTestCases = []struct {
	input    string          // input
	expected []TransitSource // expected result
}{
	{``, make([]TransitSource, 0)},
	{`[]`, make([]TransitSource, 0)},
	{`[{ "Blob":"whatever" , "Uri": "http://whatever.com/asdf?a=1&b=2"}]`, []TransitSource{TransitSource{"whatever", "http://whatever.com/asdf?a=1&b=2"}}},
	{`[{ "Blob":"s1" , "Uri": "abc"},{ "Blob":"s2" , "Uri": "abc2"}]`, []TransitSource{TransitSource{"s1", "abc"}, TransitSource{"s2", "abc2"}}},
	{`[{ "malformed json":]`, nil},
	{`[{ "Blo":"malformed name" , "Uri": "abc"}]`, nil},
}

func TestGetSources(t *testing.T) {
	for _, tc := range getSourcesTestCases {
		os.Setenv(EnvNameBilbao, tc.input)
		result := GetSources()
		if !areEqual(result, tc.expected) {
			t.Errorf("GetSources(%v): expected %v, actual %v", tc.input, tc.expected, result)
		}
	}
}

func areEqual(a, b []TransitSource) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestParse(t *testing.T) {
}

func TestLines(t *testing.T) {
}

func TestStops(t *testing.T) {
}
