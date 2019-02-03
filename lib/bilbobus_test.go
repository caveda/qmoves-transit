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
	{`[{ "Path":"whatever" , "Uri": "http://whatever.com/asdf?a=1&b=2", "Id":"whatever1"}]`, []TransitSource{TransitSource{"whatever", "http://whatever.com/asdf?a=1&b=2", "whatever1"}}},
	{`[{ "Path":"s1" , "Uri": "abc", "Id":"Id1"},{ "Path":"s2" , "Uri": "abc2", "Id":"Id2"}]`, []TransitSource{TransitSource{"s1", "abc", "Id1"}, TransitSource{"s2", "abc2", "Id2"}}},
	{`[{ "malformed json":]`, nil},
	{`[{ "Path":"malformed name" , "Uri": "abc", "Idd":"s"}, ""]`, nil},
}

func TestGetSources(t *testing.T) {
	var bilboBus Bilbobus
	for _, tc := range getSourcesTestCases {
		os.Setenv(EnvNameBilbao, tc.input)
		result := bilboBus.GetSources()
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

