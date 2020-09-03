package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/sethvargo/go-password/password"

	"google.golang.org/protobuf/proto"
)

func makeXID() string {
	a2, _ := password.Generate(2, 0, 0, true, false)
	d2, _ := password.Generate(2, 2, 0, false, false)
	a2b, _ := password.Generate(2, 0, 0, true, false)
	d2b, _ := password.Generate(2, 2, 0, false, false)

	return a2 + d2 + a2b + d2b
}

// TestPapers is the format of the output from testing
// each file is one school in this json format
type TestPapers struct {
	Region string `json:"region"`
	// Exams  []Exams `json:"exams"`
	School string `json:"school"`
}

//School is a School for output
type School struct {
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Town   string `json:"town,omitempty"`
	Active bool   `json:"active,omitempty"`
	// Region string `json:"region,omitempty"`
}

//Region is a thing that holds schools
type Region struct {
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Active bool   `json:"active,omitempty"`
	// State   string    `json:"state,omitempty"`
	Schools []*School `json:"schools,omitempty"`
}

// State is a thing with regions
type State struct {
	Name    string    `json:"name,omitempty"`
	Regions []*Region `json:"regions,omitempty"`
}

func (s *State) getRegion(sid string) (*Region, error) {
	for _, val := range s.Regions {
		if val.ID == sid {
			return val, nil
		}
	}
	return nil, errors.New("not found")
}

func (s *State) getSchool(sid string) (*School, error) {
	for _, rgn := range s.Regions {
		scol, err := rgn.getSchool(sid)
		if err == nil {
			return scol, nil
		}
	}
	return nil, errors.New("not found")
}

func (r *Region) getSchool(sid string) (*School, error) {
	for _, val := range r.Schools {
		if val.ID == sid {
			return val, nil
		}
	}
	return nil, errors.New("not found")
}

func main() {
	// urlbase := "http://localhost:8080"
	// rulesLocation := "http://www.cteresource.org/fbla/hs_competitive_events/index.html"
	schoolFilesLoc := "FBLA_2019/exams/regional/test_papers/*.json"
	// examTime := 50
	// contestInfo := "exams/regional/contest_info/contest_info.json"
	schoolFiles, _ := filepath.Glob(schoolFilesLoc)
	// loadExams := true

	// var state State

	// state.Id = "VA"
	// state.Name = "Virginia"

	// var csvFile strings.Builder

	// eventout := OutEvent{Name: "2019 Regionals"}

	var town string

	var regionsList []string

	// var schoolsCSV []*SchoolCsv

	var schoolEvents []string
	var regionEvents []string
	var schoolEventsUnmarshaled []*SchoolEvent
	var regionEventsUnmarshaled []*RegionEvent

	// var evts HierarchyEvents
	var regionID string
	// for index, element := range somelist {do something for each}
	for _, path := range schoolFiles {
		// path := schoolFiles[i]
		town = ""
		// fmt.Printf(path)
		if strings.Contains(path, "(") {
			s := strings.Index(path, "(")
			e := strings.LastIndex(path, ")")
			town = path[s+1 : e]
			// fmt.Printf(town)
		}
		jsonFile, err := os.Open(path)
		if err != nil {
			fmt.Println(err)
		}
		// fmt.Println("Successfully Opened " + path)
		// defer the closing of our jsonFile so that we can parse it later on
		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)
		var testPapers TestPapers

		json.Unmarshal(byteValue, &testPapers)
		// outpt := "region: " + testPapers.Region + " school: " + testPapers.School

		regionFound := false

		for _, val := range regionsList {
			if testPapers.Region == val {
				regionFound = true
				break
			}
		}
		if regionFound == false {
			regionsList = append(regionsList, testPapers.Region)
			fmt.Println(regionsList)
			regionID = makeXID()
			evt := &RegionEvent{Id: regionID,
				Registered: &Registered{Name: testPapers.Region,
					Parent: "VA FBLA"}}
			val, _ := proto.Marshal(evt)
			regionEvents = append(regionEvents, string(val))
			regionEventsUnmarshaled = append(regionEventsUnmarshaled, evt)
		}
		schoolID := makeXID()
		sr := &SchoolEvent{Id: schoolID,
			Registered: &Registered{
				Name:   testPapers.School,
				Parent: regionID}}
		// var reg &SchoolEvent_Registered
		// reg.Registered.Name = testPapers.School
		// reg.Registered.Parent = regionID
		// sr.Action = &reg

		val, _ := proto.Marshal(sr)
		schoolEvents = append(schoolEvents, string(val))
		schoolEventsUnmarshaled = append(schoolEventsUnmarshaled, sr)

		if len(town) > 0 {
			sr := &SchoolEvent{
				Id:     schoolID,
				Placed: &Placed{Town: town}}
			val, _ := proto.Marshal(sr)
			schoolEvents = append(schoolEvents, string(val))
			schoolEventsUnmarshaled = append(schoolEventsUnmarshaled, sr)
		}
	}
	jsonRegion, _ := json.MarshalIndent(regionEventsUnmarshaled, "", "\t")
	ioutil.WriteFile("regionEvents.json", jsonRegion, 0644)
	jsonSchool, _ := json.MarshalIndent(schoolEventsUnmarshaled, "", "\t")
	ioutil.WriteFile("schoolEvents.json", jsonSchool, 0644)

	virginia := State{Name: "Virginia"}

	for _, buf := range regionEvents {
		var event RegionEvent
		proto.Unmarshal([]byte(buf), &event)
		regionID := event.Id
		region, err := virginia.getRegion(regionID)
		if err != nil {
			region = &Region{ID: regionID}
		}
		if event.Registered != nil {
			region.Name = event.Registered.Name
			virginia.Regions = append(virginia.Regions, region)
			// region.State = event.Registered.Parent
		}
		if event.Activated != nil {
			region.Active = true
		}
		if event.Deactivated != nil {
			region.Active = false
		}
		if event.Renamed != nil {
			region.Name = event.Renamed.Name
		}
	}

	for _, buf := range schoolEvents {
		var event SchoolEvent
		proto.Unmarshal([]byte(buf), &event)
		schoolID := event.Id
		school, err := virginia.getSchool(schoolID)
		if err != nil {
			school = &School{ID: schoolID}
		}
		if event.Registered != nil {
			school.Name = event.Registered.Name
			// school.Region = event.Registered.Parent
			region, _ := virginia.getRegion(event.Registered.Parent)
			region.Schools = append(region.Schools, school)
		}
		if event.Activated != nil {
			school.Active = true
		}
		if event.Deactivated != nil {
			school.Active = false
		}
		if event.Renamed != nil {
			school.Name = event.Renamed.Name
		}
		if event.Placed != nil {
			school.Town = event.Placed.Town
		}

		fmt.Println(buf)

	}
	// decoding the events into a hierarchy!!!
	file, _ := json.MarshalIndent(virginia, "", "\t")
	ioutil.WriteFile("schools.json", file, 0644)
}
