package transit

import (
	"testing"
			"log"
	)


var remediationLinesTest = []Line{
	 Line {"I01", "01", 1, "Moon - Pluto URANO", "FORWARD", nil, nil, new(bool)},
	 Line {"I22", "22", 22, "PERSEI - GLIESE/PEGASI", "FORWARD", nil, nil, new(bool)},
	 Line {"IK4", "K4", 90881, "Tauri - Herculis", "FORWARD", nil, nil, new(bool)},
	 Line {"V01", "01", 1, "Pluto URANO - Moon", "BACKWARD", nil, nil, new(bool)},
	 Line {"V22", "22", 22, "PEGASI/GLIESE - PERSEI", "BACKWARD", nil, nil, new(bool)},
	 Line {"VK4", "K4", 90881, "Herculis - Tauris", "BACKWARD", nil, nil, new(bool)},
	 Line {"I29", "29", 29, "Orionis - Trianguli", "FORWARD", nil, nil, new(bool)},
	 Line {"V99", "99", 99, "Serpentis - Pegasi", "BACKWARD", nil, nil, new(bool)},
}


var remediateLineNameTestCases = []struct {
	id string // agencyId of line to modify
	expectedForwardName string // expected name for forward direction
	expectedBackwardName string // expected name for backwards direction
	forwardIndexChanged int // expected index element changed for forward direction
	backwardIndexChanged int // expected index element changed for backwards direction
	expectedError bool // shall return an error?
}{
	{"01", "MOON - URANUS", "URANUS - MOON",0, 3, false},
	{"22", "Arae Persei - Gliese/Pegasi", "Gliese/Pegasi - Arae Persei",1, 4, false},
	{"K4", "Ursae Majoris - Herculis", "Herculis - Ursae Majoris",2, 5, false},
	{"29", "Orionis2 - Trianguli", "Trianguli - Orionis2",0, 0, true},
	{"99", "Serpentis - Pegasis", "Serpentis - Pegasis",0, 0, true},
}

func TestRemediateLineName(t *testing.T) {
	log.Printf("---------- TestRemediateLineName ------------ ")
	for _, tc := range remediateLineNameTestCases {
		lines := make([]Line, len (remediationLinesTest))
		copy(lines, remediationLinesTest)
		err := RemediateLineName(&lines, tc.id, tc.expectedForwardName)
		// error?
		if !tc.expectedError && err != nil  {
			t.Errorf("Line %v: not expected error (%v)", tc.id, err)
		}

		if err==nil {
			// check lines values
			for j, l := range lines {
				if j == tc.forwardIndexChanged {
					if l.Name != tc.expectedForwardName {
						t.Errorf("Line %v Forward: Name expected (%v), actual (%v)", l.AgencyId, tc.expectedForwardName, l.Name)
					}
				} else if j == tc.backwardIndexChanged {
					if l.Name != tc.expectedBackwardName {
						t.Errorf("Line %v Backward: Name expected (%v), actual (%v)", l.AgencyId, tc.expectedBackwardName, l.Name)
					}
				} else {
					if l.Name != remediationLinesTest[j].Name {
						t.Errorf("Line %v. Unexpected name change. expected (%v), actual (%v)", l.AgencyId, remediationLinesTest[j].Name, l.Name)
					}
				}
			}
		}
	}
	log.Printf("------------------------------------------------ ")
}
