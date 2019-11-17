package transit

import (
	"errors"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"path"
	"strings"
)

// Constants
const GeneratedBaseNumber int = 9000
const EnvBilbobusAgencyLineWithStops = "BILBOBUS_AGENCY_LINE_WITH_STOPS"
const StopNameForwardIdentifier = "IDA"
const StopNameBackwardIdentifier = "VUELTA"
const RegexPatternStopName = `(?m).*<td headers="parada_(.*)"><span.*</span>(.*)</td>`
const RegexPatternStopId = `(?m).*codLinea=.{1,3}&amp;temporada=.{1,2}&amp;parada=(.*)\">`
const RegexPatternStopPosition = `(?m).*href="https:\/\/maps\.google\.com\/\?q=(.*),(.*)"`
const RegexPatternStopConnections = `(?mU)<td headers="correspondencias_(ida|vuelta) correspondencia_parada">(?s)(.*)<\/td>` // flag as Ungreedy
const RegexPatternStopConnectionsNumbers = `(?m).*ntido=(1|2)">\s*(\S*)\s*<\/a>`

// Globals to this file
var currentLineDirection string
var currentLinePrefixDirecton string
var currentGeneratedNumberOrdinal int = 0
var stopsCache map[string][]Stop

// StopsParser implements the signature of type Parse.
// It's responsible for adding stops to every line present in l.
func StopsParser(lines *[]Line, ts TransitSource) error {
	currentGeneratedNumberOrdinal = 0
	stopsCache = make(map[string][]Stop)

	for i, l := range *lines {
		if stopsCache[l.Id] != nil {
			(*lines)[i].Stops = stopsCache[l.Id]
		} else {
			forwardStops, backwardStops, err := fetchStopsForLine(l, ts)

			if err != nil {
				log.Printf("Error parsing stops of line %v. Error: %v ", l.AgencyId, err)
				continue
			}

			if l.Direction == DirectionForward {
				(*lines)[i].Stops = forwardStops
				stopsCache[BuildLineIdWithDirection(l.AgencyId, DirectionBackward)] = backwardStops
			} else {
				(*lines)[i].Stops = backwardStops
				stopsCache[BuildLineIdWithDirection(l.AgencyId, DirectionForward)] = forwardStops
			}
		}

		(*lines)[i].MapRoute = generateMapRoute((*lines)[i])
	}

	return nil
}

func fetchStopsForLine(l Line, ts TransitSource) (forwardStops []Stop, backwardStops []Stop, e error) {
	path, err := getLinePage(l.AgencyId, ts.Uri, ts.Path)
	if err != nil {
		log.Printf("Error getting doc of line %v. Error: %v ", l.AgencyId, err)
		return nil, nil, err
	}

	return parseLineStops(path)
}

// getLinePage downloads the document containing the line stops and
// returns the path to the document or an error if anything goes wrong.
func getLinePage(lineId string, uri string, outputDataPath string) (p string, err error) {

	log.Printf("Fetching doc with stops of line %v", lineId)

	agencyLinesUri, _ := buildStopsUri(uri, lineId)

	p = path.Join(outputDataPath, "line_stops_"+lineId+".html")
	if !UseCachedData() || !Exists(p) {
		err = Download(agencyLinesUri, p, IsFileSizeGreaterThanZero)
	}

	return p, err
}

func buildStopsUri(template, lineId string) (string, error) {
	return strings.Replace(template,
		TokenLine, lineId, 1), nil
}

// parseLineStops parses the stops file of the line a return the two lists of
// stops, one for forward direction and the other for backwards direction.
func parseLineStops(filePath string) (forwardStops []Stop, backwardStops []Stop, e error) {

	f, err := ioutil.ReadFile(filePath)
	if err != nil || len(f) == 0 {
		log.Printf("Error reading content of file %v. Error: %v ", filePath, err)
		return nil, nil, err
	}

	// extract relevant information from file
	content := string(f)
	names, err := ApplyRegexAllSubmatch(content, RegexPatternStopName)
	if err != nil || len(names) == 0 {
		log.Printf("No stop names found in content file %v. Error: %v ", filePath, err)
		return nil, nil, err
	}

	ids, err := ApplyRegexAllSubmatch(content, RegexPatternStopId)
	if err != nil || len(ids) == 0 {
		log.Printf("No stop ids found in content file %v. Error: %v ", filePath, err)
		return nil, nil, err
	}

	positions, err := ApplyRegexAllSubmatch(content, RegexPatternStopPosition)
	if err != nil || len(positions) == 0 {
		log.Printf("No stop positions found in content file %v. Error: %v ", filePath, err)
		return nil, nil, err
	}

	connectionsRaw, err := ApplyRegexAllSubmatch(content, RegexPatternStopConnections)
	if err != nil || len(positions) == 0 {
		log.Printf("No stop positions found in content file %v. Error: %v ", filePath, err)
		return nil, nil, err
	}

	if len(names) != len(ids) || len(positions) != len(names) {
		message := fmt.Sprintf("Inconsistent number of elements found by regex in stops file %v", filePath)
		log.Printf(message)
		return nil, nil, errors.New(message)
	}

	// Create the line stops from the parsed information
	var fs []Stop
	var bs []Stop
	for i, n := range names {
		d, err := getStopDirectionFromTag(n[1])
		if err != nil {
			log.Printf("Error recognizing direction of stop %v . Error: %v ", n[2], err)
			return nil, nil, err
		}

		connections := buildConnectionList(connectionsRaw[i][2])

		stop := buildStop(ids[i][1], html.UnescapeString(n[2]), connections, positions[i][1], positions[i][2])
		if d == DirectionForward {
			fs = addStopToList(fs, stop)
		} else {
			bs = addStopToList(bs, stop)
		}
	}

	log.Printf("Forward stops %v", len(fs))
	log.Printf("Backward stops %v", len(bs))
	return fs, bs, nil
}

func addStopToList(stopList []Stop, stop Stop) []Stop {

	if RemoveDuplicatedStopsInLine() {
		// Prevent duplicates
		for _, s := range stopList {
			if s.Id == stop.Id {
				log.Printf("Detected duplicated stop %v - %v", stop.Id, stop.Name)
				return stopList
			}
		}
	}

	return append(stopList, stop)
}

func getStopDirectionFromTag(tag string) (string, error) {
	tagNormalized := strings.ToUpper(strings.TrimSpace(tag))

	var d string
	if tagNormalized == StopNameForwardIdentifier {
		d = DirectionForward
	} else if tagNormalized == StopNameBackwardIdentifier {
		d = DirectionBackward
	} else {
		err := errors.New("Error: tag %v do not match any known direction: " + tagNormalized)
		return "", err
	}

	return d, nil
}

func generateMapRoute(l Line) []Coordinates {
	var route []Coordinates
	for _, s := range l.Stops {
		route = append(route, s.Location)
	}
	return route
}

func buildStop(id, name, connections, lat, long string) Stop {
	stop := Stop{id, name, connections, Timetable{"", "", "", "", ""}, Coordinates{lat, long}}
	return stop
}

func buildConnectionList(connectionsRaw string) string {
	matches, err := ApplyRegexAllSubmatch(connectionsRaw, RegexPatternStopConnectionsNumbers)
	if err != nil {
		return ""
	}

	var connections string
	for i, m := range matches {
		if !IsIgnored(m[2]) {
			connections += buildConnectionCode(m[1], m[2])
			if i < len(matches)-1 {
				connections += " "
			}
		}
	}
	return connections
}

func buildConnectionCode(direction, id string) string {
	if direction == "1" {
		return DirectionForwardShortPrefix + id
	}
	return DirectionBackwardShortPrefix + id
}

