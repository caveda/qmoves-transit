package transit

// Bilbobus is a parser of transit information of Bilbao bus agency.
type Bilbobus struct {
	lines []Line
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
