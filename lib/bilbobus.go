package transit

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"strings"
	"time"
	"strconv"
)

// Bilbobus is a parser of transit information of Bilbao bus agency.
type Bilbobus struct {
	data TransitData
}

// Constants
const EnvNameBilbao string = "BILBAO_TRANSIT"
const separator string = "-"
const gtsfStopsFileName string = "stops.txt"

var version Version = Version{"1", strconv.FormatInt(time.Now().Unix(), 10)}

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
		if s.Id == SourceDayLines {
			err = parser(&p.data.dayLines, s)
		} else if s.Id == SourceNightLines {
			err = parser(&p.data.nightLines, s)
		} else {
			err = parser(&p.data.lines, s)
		}
		if err != nil {
			log.Printf("Error while processing source %v from path %v. Error: %v ", s.Id, s.Path, err)
			continue
		}
	}
	return nil
}

// Observer that returns the list of lines for this transit.
func (p Bilbobus) Data() TransitData {
	p.data.version = version
	return p.data
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
	case SourceDayLines:
		return LinesListParser, nil
	case SourceNightLines:
		return LinesListParser, nil
	default:
		return nil, errors.New("Unknown source id " + s.Id)
	}
}
