package transit

import (
	"sync"
	"log"
	"path"
	"io/ioutil"
	"regexp"
	"fmt"
	"errors"
	"strings"
	"os"
	"time"
)

// Constants
const EnvNameBilbobusSummerStart string = "BILBOBUS_SUMMER_START"
const EnvNameBilbobusSummerEnd string = "BILBOBUS_SUMMER_END"
const scheduleRegExTemplate string = `(?m).*<a href="horario-estimado\?codLinea=` + TokenLine + `&amp;temporada=` + TokenSeason +
	`&amp;servicio=\d{0,3}&amp;tipodia=` + TokenDay + `&amp;sentido=` + TokenDirection + `&amp;hora=.*">(.*)<`
const scheduleValidationFileRegEx string = `(?m).*<a href="horario-estimado\?codLinea=`

// Types
type JobSchedule struct {
	s  *Stop
	l  Line
	ts TransitSource
}

// ScheduleParser implements the signature of type Decorator.
// It's responsible for decorating lines with the location of the stops.
func ScheduleParser(l *[]Line, ts TransitSource) error {
	c := make(chan JobSchedule)
	go scheduleMaster(c, l, ts)
	collectSchedules(3, c)
	return nil
}

func collectSchedules(workers int, c chan JobSchedule) {
	var wg sync.WaitGroup
	poolSize := 5
	for i := 0; i < poolSize; i++ {
		wg.Add(1)
		go scheduleWorker(&wg, c)
	}
	wg.Wait()
}

func scheduleMaster(c chan JobSchedule, l *[]Line, ts TransitSource) {
	for _, line := range *l {
		for j, _ := range line.Stops {
			c <- JobSchedule{&line.Stops[j], line, ts}
		}
	}
	close(c)
}

func scheduleWorker(wg *sync.WaitGroup, c <-chan JobSchedule) {
	defer wg.Done()
	for job := range c {
		log.Printf("Processing static schedule for line %v and stop %v", job.l.Id, job.s.Id)
		fillScheduleForStop(job.s, job.l, job.ts)
	}
}

// fillScheduleForStop fetch the schedule data from ts for the given stop.
// Fills out the passed Stop structure with the information fetched.
func fillScheduleForStop(s *Stop, l Line, ts TransitSource) error {
	u, err := buildScheduleUrl(ts.Uri, l.AgencyId, s.Id)
	if err != nil {
		log.Printf("Error building the url for line %v and stop %v. Error: %v ", l.Number, s.Id, err)
		return err
	}

	p := path.Join(path.Dir(ts.Path), "sched_"+l.Id+"_"+s.Id+".html")
	if !UseCachedData() || !Exists(p) {
		Download(u, p, validateScheduleFile)
	}
	err = parseScheduleFile(p, s, l)
	return err
}

func parseScheduleFile(path string, s *Stop, l Line) error {
	f, err := ioutil.ReadFile(path) // Read all
	if err != nil || len(f) == 0 {
		log.Printf("Error opening file %v. Error: %v ", path, err)
		return err
	}

	season, err := getSeason()
	if err != nil {
		log.Printf("Error figuring out the season. Error: %v ", err)
		return err
	}

	days := []string{WeekDayTypeId, SaturdayTypeId, SundayTypeId}
	for _, day := range days {
		pattern := buildScheduleRegexPattern(l.AgencyId, season, day, ToDirectionNumber(l.Direction))
		regex, err := regexp.Compile(pattern)
		if err != nil {
			log.Printf("Error compiling schedule regex %v. Error: %v ", pattern, err)
			return err
		}
		times := regex.FindAllStringSubmatch(string(f), -1)
		if times == nil {
			message := fmt.Sprintf("No static schedule found. Line: %v. Stop: %v. Day %v. Season: %v", l.Id, s.Id, day, season)
			log.Printf(message)
			continue
		}
		for i, t := range times {

			// Last element does not have separator
			separator := ""
			if i < len(times)-1 {
				separator = ","
			}

			time := string(t[1])
			if day == WeekDayTypeId {
				s.Schedule.Weekday += time + separator
			} else if day == SaturdayTypeId {
				s.Schedule.Saturday += time + separator
			} else if day == SundayTypeId {
				s.Schedule.Sunday += time + separator
			}
		}
	}

	if len(s.Schedule.Friday)==0 {
		s.Schedule.Friday = s.Schedule.Weekday
	}

	if len(s.Schedule.MondayToThrusday)==0 {
		s.Schedule.MondayToThrusday = s.Schedule.Weekday
	}

	return nil
}

func buildScheduleUrl(template, lineNumber, stopId string) (string, error) {
	season, err := getSeason()
	if err != nil {
		return "", err
	}
	return strings.Replace(strings.Replace(strings.Replace(template,
		TokenLine, lineNumber, 1),
		TokenStop, stopId, 1),
		TokenSeason, season, 1), nil
}

func buildScheduleRegexPattern(lineNumber, season, day, direction string) string {
	return strings.Replace(strings.Replace(strings.Replace(strings.Replace(scheduleRegExTemplate,
		TokenLine, lineNumber, 1),
		TokenDay, day, 1),
		TokenSeason, season, 1),
		TokenDirection, direction, 1)
}

func getSeason() (string, error) {
	summerStart := os.Getenv(EnvNameBilbobusSummerStart)
	if len(summerStart) == 0 {
		return "", errors.New(fmt.Sprintf("Warning: Env variable %v is empty!", EnvNameBilbobusSummerStart))
	}
	summerEnd := os.Getenv(EnvNameBilbobusSummerEnd)
	if len(summerEnd) == 0 {
		return "", errors.New(fmt.Sprintf("Warning: Env variable %v is empty!", EnvNameBilbobusSummerEnd))
	}
	startTime, err := time.Parse(time.RFC3339, summerStart)
	if err != nil {
		return "", err
	}
	endTime, err := time.Parse(time.RFC3339, summerEnd)
	if err != nil {
		return "", err
	}
	season := SeasonSummer
	if time.Now().Before(startTime) || time.Now().After(endTime) {
		season = SeasonWinter
	}
	return season, nil
}


func validateScheduleFile (p string) bool {
	fi, err := os.Stat(p)
	if err != nil {
		log.Printf("validateScheduleFile: Error gettings stats file %v. Error: %v ", p, err)
		return false
	}

	// checks the size
	if fi.Size()==0 {
		log.Printf("validateScheduleFile: file %v has zero size", p)
		return false
	}

	f, err := ioutil.ReadFile(p) // Read all
	if err != nil || len(f) == 0 {
		log.Printf("validateScheduleFile: Error opening file %v. Error: %v ", p, err)
		return false
	}

	regex, err := regexp.Compile(scheduleValidationFileRegEx)
	if err != nil {
		log.Printf("validateScheduleFile: Error compiling schedule regex %v. Error: %v ", scheduleValidationFileRegEx, err)
		return false
	}
	times := regex.FindAllStringSubmatch(string(f), -1)
	if times == nil {
		log.Printf("validateScheduleFile: No static schedule found with pattern %v", scheduleValidationFileRegEx)
		return false
	}

	// Valid
	return true
}