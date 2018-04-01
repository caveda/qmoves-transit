// Package transit provides functions for fetching and processing public transit
// information.
package transit

type Coordinates struct {
	lat, long string
}

type Stop struct {
	id          string
	order       int
	name        string
	location    Coordinates
	connections []string
}

type Line struct {
	id        string
	name      string
	direction string
	stops     []Stop
}

// Parser is an interface that must be implemented per transit agency.
// Exposes Parse method to digest the raw data provided by the agency.
// Once parsed, the agency information can be queries using the rest of the methods: Lines, Stops, etc.
type Parser interface {
	Parse(dataPath string) error
	Lines() []Line
	Stops(lineId string, direction string) []Stop
}
