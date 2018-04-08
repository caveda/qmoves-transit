package transit

import (
	"encoding/json"
	"fmt"
)

// Returns the lines with the right format to be presented.
// Tipically the chosen format is json.
func Format(lines []Line) (string, error) {
	b, err := json.Marshal(lines)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return string(b), nil
}
