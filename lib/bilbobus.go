package transit

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"
)

const EnvNameBilbao string = "BILBAO_TRANSIT"

// Bilbobus is a parser of transit information of Bilbao bus agency.
type Bilbobus struct {
	lines []Line
}

func GetSources() []TransitSource {
	envData := os.Getenv(EnvNameBilbao)
	var sources []TransitSource
	dec := json.NewDecoder(strings.NewReader(envData))
	for {
		if err := dec.Decode(&sources); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
	}
	return sources
}

// Process the data files in folder dataPath
func (p Bilbobus) Parse(dataPath string) error {
	return nil
}

func (p Bilbobus) Lines() []Line {
	return p.lines
}

func (p Bilbobus) Stops(lineId string, direction string) []Stop {
	return p.lines[0].stops
}

func (p Bilbobus) digestStops(dataPath string) error {
	return nil
}

func (p Bilbobus) digestLines(dataPath string) error {
	return nil
}
