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
	"fmt"
)

// Bilbobus is a parser of transit information of Bilbao bus agency.
type Bilbobus struct {
	data TransitData
}

// Constants
const EnvNameBilbao string = "BILBAO_TRANSIT"
const separator string = "-"
const EnvMetadata string = "METADATA"

var defaultMetadataItem = MetadataItem{"1", "1", "1", "86400", "False", strconv.FormatInt(time.Now().Unix(), 10)}

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
			log.Printf("Error getting parser or source %v: %v", s.Id, e)
			continue
		}

		err = parser(&p.data.lines, s)

		if err != nil {
			log.Printf("Error while processing source %v from path %v. Error: %v ", s.Id, s.Path, err)
			continue
		}
	}

	// All sources processed. Add the list of stops
	p.data.stops, _ = extractStops(p.data.lines)
	return nil
}

// Observer that returns the list of lines for this transit.
func (p Bilbobus) Data() TransitData {
	p.data.metadata = p.BuildMetadata()
	return p.data
}

// getParser returns the proper parser for the given transitSource.
func getParser(s TransitSource) (Parse, error) {
	switch s.Id {
	case SourceLines:
		return LinesParser, nil
	case SourceStops:
		return StopsParser, nil
	case SourceSchedule:
		return ScheduleParser, nil
	default:
		return nil, errors.New("Unknown source id " + s.Id)
	}
}

// tagNightlyLines walk through the list of lines setting
// to true those considered nightly
func tagNightlyLines(t *TransitData) error {

	// Precondition
	if t.lines==nil || t.nightLines == nil || len(t.lines) <= 0 || len(t.nightLines) <= 0 {
		message := fmt.Sprintf("Either lines or nightlines has no elements")
		log.Printf(message)
		return errors.New(message)
	}

	for _, nl := range t.nightLines {
		for i, l := range t.lines {
			if nl.Number==l.Number {
				*(t.lines[i].IsNightLine) = true
			}
		}
	}
	return nil
}

// extractStops explores the supplied transit data to fill out
// the list of stops from the rest of the information contained.
func extractStops (l []Line) ([]Stop, error){
	stops := make(map[string]Stop)
	for _, line := range l {
		for _, s := range line.Stops {
			_, stopPresent := stops[s.Id]
			if !stopPresent {
				s.Connections = addLineIdToConnections(s.Connections, line.Id)
				stops[s.Id] = s
			}
		}
	}

	return toStopSlice(stops), nil
}

func addLineIdToConnections (connections string, newLineId string) string {

	if len(connections)!=0 {
		connections += " "
	}
	connections+=newLineId
	return connections
}

func toStopSlice(m map[string]Stop) []Stop {
	v := make([]Stop, len(m))
	index := 0
	for _, value := range m {
		v[index] = value
		index++
	}
	return v
}

// Read the metadata from the env var.
// Returns an object holding the Metadata information.
func (p Bilbobus) BuildMetadata() []MetadataItem {
	envData := os.Getenv(EnvMetadata)
	// Base case
	if len(envData) == 0 {
		log.Printf("Warning: Env variable %v is empty!", EnvMetadata)
		metadataDefault := make([]MetadataItem, 0)
		return append(metadataDefault, defaultMetadataItem)
	}

	var metadata []MetadataItem
	dec := json.NewDecoder(strings.NewReader(envData))
	dec.DisallowUnknownFields()
	for {
		if err := dec.Decode(&metadata); err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Error while parsing input: %v ", err)
			return nil
		}
	}

	// Insert last update timestamp to the one which has not
	for index, value := range metadata {
		if len(value.LastUpdate)==0 {
			metadata[index].LastUpdate = strconv.FormatInt(time.Now().Unix(), 10)
			break
		}
	}
	return metadata
}