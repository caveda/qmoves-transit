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
}{
	{`ARANGOITI - PLAZA BIRIBILA`, `Arangoiti - Gran Via`, DirectionForward, DirectionForwardShortPrefix},
	{`ARANGOITI - PLAZA BIRIBILA`, `Gran Via - Arangoiti`, DirectionBackward, DirectionBackwardShortPrefix},
	{`PLAZA BIRIBILA - OTXARKOAGA`, `Plaza Biribila/GV  - Otxarkoaga`, DirectionForward, DirectionForwardShortPrefix},
	{`PLAZA BIRIBILA - OTXARKOAGA`, `Otxarkoaga - Plaza Biribila/GV`, DirectionBackward, DirectionBackwardShortPrefix},
	{`SAN MAMES - ARABELLA`, `San Mames - Arabella`, DirectionForward, DirectionForwardShortPrefix},
	{`SAN MAMES - ARABELLA`, `Arabella - San Mames`, DirectionBackward, DirectionBackwardShortPrefix},
}

func TestGetLineDirection(t *testing.T) {
	for _, tc := range getLineDirectionTestCases {
		long, short := GetLineDirection(tc.lineName, tc.rawDirection)
		if long != tc.expectedLong || short != tc.expectedShort {
			t.Errorf("getLineDirection(%v,%v): expected (%v,%v), actual (%v,%v)", tc.lineName, tc.rawDirection, tc.expectedLong, tc.expectedShort, long, short)
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
