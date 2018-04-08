// Package transit provides functions for fetching and processing public transit
// information.
package transit

type TransitSource struct {
	Blob, Uri string
}

type Coordinates struct {
	lat, long string
}

type Stop struct {
	Id          string `json:"Id,omitempty"`
	Name        string `json:"Name,omitempty"`
	Connections string `json:"Connections,omitempty"`
	Schedule    string `json:"Schedule,omitempty"`
}

type Line struct {
	Id        string `json:"Id,omitempty"`
	Name      string `json:"Name,omitempty"`
	Direction string `json:"Direction,omitempty"`
	Stops     []Stop `json:"Stops,omitempty"`
}

// Parser is an interface that must be implemented per transit agency.
// Exposes Parse method to digest the raw data provided by the agency.
// Once parsed, the agency information can be queries using the rest of the methods: Lines, Stops, etc.
type Parser interface {
	Parse(dataPath string) error
	Lines() []Line
	Stops(lineId string, direction string) []Stop
}