package transit

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"strings"
)

// Bilbobus is a parser of transit information of Bilbao bus agency.
type Bilbobus struct {
	lines []Line
}

// Constants
const EnvNameBilbao string = "BILBAO_TRANSIT"
const separator string = "-"
const gtsfStopsFileName string = "stops.txt"

// Read the data sources from the env var.
// Returns the list of sources.
func (p Bilbobus) GetSources() []TransitSource {
	envData := os.Getenv(EnvNameBilbao)
	// Base case
	if len(envData) == 0 {
		log.Printf("Warning: Env variable %v is empty!", EnvNameBilbao)
		return make([]TransitSource, 0)
	}

	var sources []TransitSource
	dec := json.NewDecoder(strings.NewReader(envData))
	dec.DisallowUnknownFields()
	for {
		if err := dec.Decode(&sources); err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Error while parsing input: %v ", err)
			return nil
		}
	}
	return sources
}

// Process the data files in folder dataPath and build the data model.
func (p *Bilbobus) Digest(sources []TransitSource) error {
	var err error
	for _, s := range sources {
		parser, e := getParser(s)
		if e != nil {
			log.Printf("Error getting parser for source %v: %v", s.Id, e)
			continue
		}
		err = parser(&p.lines, s)
		if err != nil {
			log.Printf("Error while processing source %v from path %v. Error: %v ", s.Id, s.Path, err)
			continue
		}
	}
	return nil
}

// Observer that returns the list of lines for this transit.
func (p Bilbobus) Lines() []Line {
	return p.lines
}

// getParser returns the proper parser for the given transitSource.
func getParser(s TransitSource) (Parse, error) {
	switch s.Id {
	case SourceLines:
		return LinesParser, nil
	case SourceLocation:
		return LocationParser, nil
	case SourceSchedule:
		return ScheduleParser, nil
	default:
		return nil, errors.New("Unknown source id " + s.Id)
	}
}

func GetLineDirection(name string, rawDirection string) (long string, short string, err error) {
	nameBegin := strings.ToUpper(strings.TrimSpace(strings.Split(name, separator)[0]))
	directionBegin := strings.ToUpper(strings.TrimSpace(strings.Split(rawDirection, separator)[0]))
	nameEnd := strings.ToUpper(strings.TrimSpace(strings.Split(name, separator)[1]))
	directionEnd := strings.ToUpper(strings.TrimSpace(strings.Split(rawDirection, separator)[1]))
	long = DirectionBackward
	err = nil

	// Check error control
	if !(strings.Contains(directionEnd, nameEnd) || strings.Contains(directionBegin, nameBegin)) && !(strings.Contains(directionBegin, nameEnd) || strings.Contains(directionEnd, nameBegin)) {
		err = errors.New("Error: Name and direction do not match. Direction: " + directionBegin + "," + directionEnd + " .Name: " + nameBegin + "," + nameEnd)
		return "", "", err
	}

	if strings.Contains(directionBegin, nameBegin) {
		long = DirectionForward
	}

	return long, ToDirectionPrefix(long), err
}
