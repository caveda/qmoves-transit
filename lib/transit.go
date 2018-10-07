// Package transit specifies the model (types, const, data structures,...) of
// the library.
package transit

import (
	"os"
	"strconv"
)

// Consts
const SourceLines string = "Lines"
const SourceSchedule string = "Schedule"
const SourceLocation string = "Location"
const SourceDayLines string = "DayLinesList"
const SourceNightLines string = "NightLinesList"
const DirectionForward string = "FORWARD"
const DirectionBackward string = "BACKWARD"
const DirectionForwardShortPrefix string = "I"
const DirectionBackwardShortPrefix string = "V"

const TokenLine string = "<LINEID>"
const TokenSeason string = "<SEASON>"
const TokenStop string = "<STOPID>"
const TokenDirection string = "<DIRECTIONID>"
const TokenDay string = "<DAYTYPE>"
const SeasonWinter string = "IV"
const SeasonSummer string = "VE"
const WeekDayTypeId string = "1"
const SaturdayTypeId string = "2"
const SundayTypeId string = "3"
const DirectionForwardNumber string = "1"
const DirectionBackwardNumber string = "2"
const EnvNameReuseLocalData string = "REUSE_TRANSIT_LOCAL_FILES"

// Globals
var Directions = [2]string{DirectionForward, DirectionBackward}
var DirectionsPrefixes = [2]string{DirectionForwardShortPrefix, DirectionBackwardShortPrefix}

// Bilbobus is a parser of transit information of Bilbao bus agency.
type TransitData struct {
	version    	Metadata
	lines      	[]Line
	dayLines   	[]Line
	nightLines 	[]Line
	stops		[]Stop
}

// Metadata contains meta-information about the data
// provided, such as: last time data was updated, version, etc.
type Metadata struct {
	Version    string `json:"Version,omitempty"`
	LastUpdate string `json:"LastUpdate,omitempty"`
}

// TransitSource tells what data a source has to have.
type TransitSource struct {
	Path, Uri, Id string
}

// Location data (typically of a stop).
type Coordinates struct {
	Lat  string `json:"Lat,omitempty"`
	Long string `json:"Long,omitempty"`
}

// Timetable stores the schedule per type of day.
type Timetable struct {
	Weekday  string `json:"Weekday,omitempty"`
	MondayToThrusday string `json:"MondayToThrusday,omitempty"`
	Friday   string `json:"Friday,omitempty"`
	Saturday string `json:"Saturday,omitempty"`
	Sunday   string `json:"Sunday,omitempty"`
}

// Stop keeps the information of a (bus, metro,...) stop.
type Stop struct {
	Id          string      `json:"Id,omitempty"`
	Name        string      `json:"Name,omitempty"`
	Connections string      `json:"Connections,omitempty"`
	Schedule    Timetable   `json:"Schedule,omitempty"`
	Location    Coordinates `json:"Location,omitempty"`
}

// Line represents a line of transport mean. Consists of
// some data and the list of stops for the route of the line.
type Line struct {
	Id        string        `json:"Id,omitempty"`
	AgencyId  string		`json:"AgencyId,omitempty"`
	Number    int       	`json:"Number,omitempty"`
	Name      string        `json:"Name,omitempty"`
	Direction string        `json:"Direction,omitempty"`
	Stops     []Stop        `json:"Stops,omitempty"`
	MapRoute  []Coordinates `json:"MapRoute,omitempty"`
	IsNightLine *bool		`json:"IsNightLine,omitempty"`
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

// ToDirectionNumber returns the identifier that matches
// the given direction string (either DirectionBackward or DirectionForward )
func ToDirectionNumber(direction string) string {
	id := DirectionForwardNumber
	if direction == DirectionBackward {
		id = DirectionBackwardNumber
	}
	return id
}

// ToDirectionPrefix returns the prefix that matches
// the given direction string (either DirectionBackward or DirectionForward )
func ToDirectionPrefix(direction string) string {
	p := DirectionForwardShortPrefix
	if direction == DirectionBackward {
		p = DirectionBackwardShortPrefix
	}
	return p
}

// UseCachedData returns True if the locally cached data must be used
// as data source for transit information.
func UseCachedData() bool {
	result := false
	value := os.Getenv(EnvNameReuseLocalData)
	b, err := strconv.ParseBool(value)
	if err == nil {
		result = b
	}
	return result
}
