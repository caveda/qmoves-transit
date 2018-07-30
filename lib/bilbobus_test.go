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

var getLineDirectionTestCases = []struct {
	lineName, rawDirection      string // input
	expectedLong, expectedShort string // expected result
	expectedError               bool
}{
	{`ARANGOITI - PLAZA BIRIBILA`, `Arangoiti - Gran Via`, DirectionForward, DirectionForwardShortPrefix, false},
	{`ARANGOITI - PLAZA BIRIBILA`, `Gran Via - Arangoiti`, DirectionBackward, DirectionBackwardShortPrefix, false},
	{`PLAZA BIRIBILA - OTXARKOAGA`, `Plaza Biribila/GV  - Otxarkoaga`, DirectionForward, DirectionForwardShortPrefix, false},
	{`PLAZA BIRIBILA - OTXARKOAGA`, `Otxarkoaga - Plaza Biribila/GV`, DirectionBackward, DirectionBackwardShortPrefix, false},
	{`SAN MAMES - ARABELLA`, `San Mames - Arabella`, DirectionForward, DirectionForwardShortPrefix, false},
	{`SAN MAMES - ARABELLA`, `Arabella - San Mames`, DirectionBackward, DirectionBackwardShortPrefix, false},
	{`LA PE�A - PLAZA BIRIBILA`, `Zamakola168 - Ayala`, "", "", true},
	{`LA PE�A - PLAZA BIRIBILA`, `Ayala - Zamakola168`, "", "", true},
}

func TestGetLineDirection(t *testing.T) {
	for _, tc := range getLineDirectionTestCases {
		long, short, err := GetLineDirection(tc.lineName, tc.rawDirection)
		if long != tc.expectedLong || short != tc.expectedShort || err != nil != tc.expectedError {
			t.Errorf("getLineDirection(%v,%v): expected (%v,%v,%v), actual (%v,%v,%v)", tc.lineName, tc.rawDirection, tc.expectedLong, tc.expectedShort, tc.expectedError, long, short, err)
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
