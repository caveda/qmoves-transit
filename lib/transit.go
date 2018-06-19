// Package transit specifies the model (types, const, data structures,...) of
// the library.
package transit

// Consts
const SourceLines string = "Lines"
const SourceSchedule string = "Schedule"
const SourceLocation string = "Location"

// Types

// TransitSource tells what data a source has to have.
type TransitSource struct {
	Path, Uri, Id string
}

type Coordinates struct {
	Lat  string `json:"Lat,omitempty"`
	Long string `json:"Long,omitempty"`
}

type Stop struct {
	Id          string      `json:"Id,omitempty"`
	Name        string      `json:"Name,omitempty"`
	Connections string      `json:"Connections,omitempty"`
	Schedule    string      `json:"Schedule,omitempty"`
	Location    Coordinates `json:"Location,omitempty"`
}

type Line struct {
	Id        string `json:"Id,omitempty"`
	Name      string `json:"Name,omitempty"`
	Direction string `json:"Direction,omitempty"`
	Stops     []Stop `json:"Stops,omitempty"`
}

// Parse is a type of function that receives a list of Lines and adds
// new information to that list. For example: a Decorator function might
// add the location of each stop in the provided list of lines.
type Parse func(*[]Line, TransitSource) error

// Parser is an interface that must be implemented per transit agency.
// Exposes Parse method to digest the raw data provided by the agency.
// Once parsed, the agency information can be queries using the rest of the methods: Lines, Stops, etc.
type Parser interface {
	Digest(dataPath string) error
	Lines() []Line
	Stops(lineId string, direction string) []Stop
}

// Presenter is an interface implemented by formatter classes.
// Exposes methods to transform lines into the presenting format, e.g. JSON
type Presenter interface {
	Format(l Line) (string, error)
	FormatList(l []Line) (string, error)
}
