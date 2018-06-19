package transit

import (
	"encoding/json"
	"fmt"
)

// Bilbobus is a parser of transit information of Bilbao bus agency.
type JsonPresenter struct {
}

// Returns line with the right format to be presented.
// Tipically the chosen format is json.
func (p JsonPresenter) Format(l Line) (string, error) {
	b, err := json.MarshalIndent(l, "", "    ")
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return string(b), nil
}

// Returns the array of lines with the right format to be presented.
// Tipically the chosen format is json.
func (p JsonPresenter) FormatList(l []Line) (string, error) {
	b, err := json.MarshalIndent(l, "", "    ")
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return string(b), nil
}
