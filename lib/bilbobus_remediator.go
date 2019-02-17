package transit

import (
	"fmt"
	"log"
	"strings"
	"errors"
)

const RemediationTag = "REMEDIATE"
const agencyNameSeparator string = "-"

// RemediateLineName replace the detected line name in data by the
// expected. The expected name shall correspond to direction FORWARD.
func  RemediateLineName(lines *[]Line, agencyLineId string, expected string) error {
	message := fmt.Sprintf(RemediationTag + ": Changing line %v name. Expected: %v", agencyLineId, expected)
	log.Printf(message)

	forwardDone := false
	backwardDone := false
	for i, l := range *lines {
		if l.AgencyId == agencyLineId {
			if l.Direction==DirectionForward {
				(*lines)[i].Name = expected
				forwardDone = true
			} else {
				backwardName, err := ReverseLineName(expected)
				if err!=nil {
					return err
				}
				(*lines)[i].Name = backwardName
				backwardDone = true
			}
		}

		// Have we finished?
		if forwardDone && backwardDone {
			log.Printf(RemediationTag + ": Remediation done successfully")
			return nil
		}
	}

	// If we reach this point. Remediation is not done.
	err := errors.New(RemediationTag + ": Name of line " + agencyLineId + " not changed to " + expected)
	return err
}

// ReverseLineName takes a line name formatted as "origin - destination" and
// returns "destination - origin"
func ReverseLineName (name string) (string, error) {

	nameParts := strings.Split(name, agencyNameSeparator)
	if len(nameParts) != 2 {
		err := errors.New(RemediationTag + ": Name " + name + " can not be splitted in two parts using separator %v" + agencyNameSeparator)
		return "", err
	}
	origin := strings.TrimSpace(nameParts[0])
	destination := strings.TrimSpace(nameParts[1])
	return destination + " - " + origin, nil
}