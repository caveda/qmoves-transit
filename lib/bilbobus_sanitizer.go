package transit

import (
				"github.com/pkg/errors"
	"strings"
	)


// Structs
type AgencyLine struct {
	Id        string
	Name      string
}

// CheckConsistency verifies that the output data is consistent
// with the information published by the transit agency
func CheckConsistency(td TransitData) (report string, err error) {

	var linesReport string
	if linesReport, err = checkLines(td.lines); err != nil {
		return  linesReport, err
	}

	report += linesReport
	return report, nil
}


// checkLines verifies that the lines information is complete.
// Returns a string with the report of verification and an error
// if any found issue couldn't be solved.
func checkLines(lines []Line) (report string, err error) {

	var str strings.Builder
	str.WriteString("\n------ Lines check -------")

	for _, l := range lines {
		str.WriteString("\nLine " + l.Id + ": ")

		lineError := false
		if lineError = len(l.Stops) == 0; lineError {
			str.WriteString("No stops.")
		} else if lineError = l.Name == ""; lineError {
			str.WriteString("No name.")
		} else if lineError = l.Number == 0; lineError {
			str.WriteString("No number.")
		} else if lineError = l.AgencyId == ""; lineError {
			str.WriteString("No agencyId.")
		} else if lineError = len(l.MapRoute) == 0; lineError {
			str.WriteString("No map route.")
		} else {
			var stopsReport string
			if stopsReport, lineError = checkStops(l.Stops); lineError {
				str.WriteString(stopsReport + ". Line deleted")
			} else {
				str.WriteString("Ok")
			}
		}

		if lineError {
			return str.String(), errors.New(str.String())
		}
	}
	return str.String(), nil
}


// checkStops verifies that each stop in the list is complete.
// Returns a string with the report of verification and an error
// if any found issue couldn't be solved.
func checkStops (stops []Stop) (report string, error bool) {

	var str strings.Builder
	for _, s := range stops {
		if error = s.Location.Lat=="" || s.Location.Long==""; error {
			str.WriteString("No location for stop " + s.Id)
			break
		}
		if error = s.Schedule.MondayToThrusday=="" && s.Schedule.Friday=="" && s.Schedule.Weekday=="" &&
			s.Schedule.Saturday=="" && s.Schedule.Sunday==""; error{
			str.WriteString("No schedule for stop " + s.Id)
			break
		}
		if error = s.Name==""; error {
			str.WriteString("No name for stop " + s.Id)
			break
		}
	}
	return str.String(),error
}


// deleteLine deletes an element from the list of lines.
//func deleteLine (lines []Line, pos int) []Line {
//
//}


// RemediateEnabled returns True if the errors
// found by sanitizer must be remediated as detected.
//func RemediateEnabled () bool {
//	result := false
//	value := os.Getenv(EnvNameReuseLocalData)
//	b, err := strconv.ParseBool(value)
//	if err == nil {
//		result = b
//	}
//	return result
//}