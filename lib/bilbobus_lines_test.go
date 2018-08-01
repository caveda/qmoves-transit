package transit

import (
	"testing"
)

var getLineDirectionTestCases = []struct {
	lineName, rawDirection      string // input
	expectedLong, expectedShort string // expected result
	expectedError               bool
}{
	{`EQUIDEM - ITA SENTIO`, `Equidem - Peculiarem in`, DirectionForward, DirectionForwardShortPrefix, false},
	{`EQUIDEM - ITA SENTIO`, `Peculiarem in - Equidem`, DirectionBackward, DirectionBackwardShortPrefix, false},
	{`STUDIIS CAUSAM - EORUM`, `Studiis Causam/GV  - Eorum`, DirectionForward, DirectionForwardShortPrefix, false},
	{`STUDIIS CAUSAM - EORUM`, `Eorum - Studiis Causam/GV`, DirectionBackward, DirectionBackwardShortPrefix, false},
	{`QUI DIFFICULTATIBUS - ESSE`, `Qui Difficultatibus - Esse`, DirectionForward, DirectionForwardShortPrefix, false},
	{`QUI DIFFICULTATIBUS - ESSE`, `Esse - Qui Difficultatibus`, DirectionBackward, DirectionBackwardShortPrefix, false},
	{`VICTIS UTILITATEMÑ - GRATIAE PLACENDI`, `Iuvandi443 - Praetulerint`, "", "", true},
	{`VICTIS UTILITATEMÑ - GRATIAE PLACENDI`, `Praetulerint - Iuvandi443`, "", "", true},
	{`IDQUE - IAM`, `Et In - Iam`, DirectionForward, DirectionForwardShortPrefix, false},
	{`IDQUE - IAM`, `Iam - Et In`, DirectionBackward, DirectionBackwardShortPrefix, false},
	{`ALIIS OPERIBUS - IPSE`, `Aliis Operibus - Ipse`, DirectionForward, DirectionForwardShortPrefix, false},
	{`ALIIS OPERIBUS - IPSE`, `Ipse - Aliis Operibus`, DirectionBackward, DirectionBackwardShortPrefix, false},
	{`PROFITEOR - MIRARI/LIVIUM`, `Auctorem celeberrimum - Profiteor  - Mirari/Livium`, "", "", true},
	{`PROFITEOR - MIRARI/LIVIUM`, `Mirari/Auctorem celeberrimum - Profiteor`, DirectionBackward, DirectionBackwardShortPrefix, false},
}

func TestGetLineDirection(t *testing.T) {
	for _, tc := range getLineDirectionTestCases {
		long, short, err := GetLineDirection(tc.lineName, tc.rawDirection)
		if long != tc.expectedLong || short != tc.expectedShort || err != nil != tc.expectedError {
			t.Errorf("getLineDirection(%v,%v): expected (%v,%v,%v), actual (%v,%v,%v)", tc.lineName, tc.rawDirection, tc.expectedLong, tc.expectedShort, tc.expectedError, long, short, err != nil)
		}
	}
}
