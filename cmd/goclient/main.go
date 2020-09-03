package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/sethvargo/go-password/password"
	"google.golang.org/protobuf/proto"
)

// TestPapers is the format of the output from testing
// each file is one school in this json format
type TestPapers struct {
	Region string  `json:"region"`
	Exams  []Exams `json:"exams"`
	School string  `json:"school"`
}

// Responses is really wrong and must be fixed eventually
// it's the student's response for a particular question by number
// usually a, b, c, d, etc.
type Responses struct {
	Num100 string `json:"100"`
	Num098 string `json:"098"`
	Num099 string `json:"099"`
	Num090 string `json:"090"`
	Num091 string `json:"091"`
	Num092 string `json:"092"`
	Num093 string `json:"093"`
	Num094 string `json:"094"`
	Num095 string `json:"095"`
	Num096 string `json:"096"`
	Num097 string `json:"097"`
	Num010 string `json:"010"`
	Num011 string `json:"011"`
	Num012 string `json:"012"`
	Num013 string `json:"013"`
	Num014 string `json:"014"`
	Num015 string `json:"015"`
	Num016 string `json:"016"`
	Num017 string `json:"017"`
	Num018 string `json:"018"`
	Num019 string `json:"019"`
	Num025 string `json:"025"`
	Num024 string `json:"024"`
	Num027 string `json:"027"`
	Num026 string `json:"026"`
	Num021 string `json:"021"`
	Num020 string `json:"020"`
	Num023 string `json:"023"`
	Num022 string `json:"022"`
	Num029 string `json:"029"`
	Num028 string `json:"028"`
	Num038 string `json:"038"`
	Num039 string `json:"039"`
	Num032 string `json:"032"`
	Num033 string `json:"033"`
	Num030 string `json:"030"`
	Num031 string `json:"031"`
	Num036 string `json:"036"`
	Num037 string `json:"037"`
	Num034 string `json:"034"`
	Num035 string `json:"035"`
	Num049 string `json:"049"`
	Num048 string `json:"048"`
	Num047 string `json:"047"`
	Num046 string `json:"046"`
	Num044 string `json:"044"`
	Num043 string `json:"043"`
	Num042 string `json:"042"`
	Num041 string `json:"041"`
	Num040 string `json:"040"`
	Num058 string `json:"058"`
	Num059 string `json:"059"`
	Num054 string `json:"054"`
	Num055 string `json:"055"`
	Num056 string `json:"056"`
	Num057 string `json:"057"`
	Num050 string `json:"050"`
	Num051 string `json:"051"`
	Num052 string `json:"052"`
	Num053 string `json:"053"`
	Num061 string `json:"061"`
	Num060 string `json:"060"`
	Num063 string `json:"063"`
	Num062 string `json:"062"`
	Num065 string `json:"065"`
	Num064 string `json:"064"`
	Num067 string `json:"067"`
	Num066 string `json:"066"`
	Num069 string `json:"069"`
	Num068 string `json:"068"`
	Num076 string `json:"076"`
	Num077 string `json:"077"`
	Num074 string `json:"074"`
	Num075 string `json:"075"`
	Num072 string `json:"072"`
	Num073 string `json:"073"`
	Num070 string `json:"070"`
	Num071 string `json:"071"`
	Num045 string `json:"045"`
	Num078 string `json:"078"`
	Num079 string `json:"079"`
	Num089 string `json:"089"`
	Num088 string `json:"088"`
	Num083 string `json:"083"`
	Num082 string `json:"082"`
	Num081 string `json:"081"`
	Num080 string `json:"080"`
	Num087 string `json:"087"`
	Num086 string `json:"086"`
	Num085 string `json:"085"`
	Num084 string `json:"084"`
	Num003 string `json:"003"`
	Num002 string `json:"002"`
	Num001 string `json:"001"`
	Num007 string `json:"007"`
	Num006 string `json:"006"`
	Num005 string `json:"005"`
	Num004 string `json:"004"`
	Num009 string `json:"009"`
	Num008 string `json:"008"`
}

// Exams is metadata for a test paper
type Exams struct {
	Status string `json:"status"`
	// for the moment, we don't need the responses in this struct
	// Responses       Responses   `json:"responses"`
	Students        []string    `json:"students" datastore:"flatten"`
	StudentMetadata interface{} `json:"student_metadata"`
	TimedTest       bool        `json:"timed_test"`
	TeamID          interface{} `json:"team_id"`
	Score           int         `json:"score"`
	ExamID          string      `json:"exam_id"`
	Sid             string      `json:"sid"`
	TimeRemaining   float64     `json:"time_remaining"`
	Password        string      `json:"password"`
}

// Student holds student name and gradelevel
type Student struct {
	Name       string `json:"name"`
	GradeLevel string `json:"grade_level,omitempty"`
}

// ExamInfo is the metadata about the
// exams in an event
type ExamInfo struct {
	Base  Base   `json:"base"`
	Exams []Exam `json:"exams"`
}

// Base is for competition-wide data
type Base struct {
	Manual    string `json:"manual"`
	PngOffset string `json:"pngOffset"`
	Level     string `json:"level"`
}

// Exam is registration info about a particular
// exam in a competition
type Exam struct {
	CollaborationMax int    `json:"collaborationMax"`
	Title            string `json:"title"`
	MinGradeLevel    string `json:"minGradeLevel"`
	Registration     string `json:"registration"`
	MaxTeams         int    `json:"maxTeams"`
	ID               string `json:"id"`
	MaxGradeLevel    string `json:"maxGradeLevel"`
	Filename         string `json:"filename"`
	PapersPerTeam    int    `json:"papersPerTeam,omitempty"`
	NeedGradeLevel   bool   `json:"needGradeLevel,omitempty"`
}

// Event is metadata about the competition event
type Event struct {
	State        string `json:"state"`
	Level        string `json:"level"`
	Organization string `json:"organization"`
	Year         string `json:"year"`
	Title        string `json:"title"`
	AdminEmail   string `json:"admin_email"`
	RulesURL     string `json:"rules_url"`
}

// Examination has the items for an exam
type Examination struct {
	Items []Item `json:"items"`
	Title string `json:"title"`
	ID    string `json:"id"`
	Time  int    `json:"time"`
}

// Choice holds a response to a question. e.g., four
// of these for most multiple-choice questions
type Choice struct {
	ID       string `json:"id"`
	Response string `json:"response"`
}

//Item holds the data for a standard question item
type Item struct {
	ItemNumber    string   `json:"item_number"`
	CorrectAnswer string   `json:"correct_answer"`
	Question      string   `json:"question"`
	Rubric        string   `json:"rubric"`
	QuestionType  string   `json:"question_type"`
	Choices       []Choice `json:"choices"`
}

//OutEvent is an Event for outputting JSON hierarchy
type OutEvent struct {
	Name    string       `json:"name"`
	Regions []*OutRegion `json:"regions"`
}

func (e *OutEvent) getRegion(region string) (*OutRegion, error) {
	ErrNotFound := errors.New("not Found")
	for _, val := range e.Regions {
		if val.Name == region {
			return val, nil
		}
	}
	return nil, ErrNotFound
}

func (e *OutEvent) getAllXIDs(known []string) []string {
	for _, region := range e.Regions {
		known = append(known, region.Xid)
		var schoolIds []string
		for _, school := range region.Schools {
			schoolIds = append(schoolIds, school.Xid)
		}
		known = append(known, schoolIds...)
	}
	return known
}

func (e *OutEvent) makeXIDS() {
	for _, region := range e.Regions {
		region.makeXIDS()
	}
}

func makeXID() string {
	a2, _ := password.Generate(2, 0, 0, true, false)
	d2, _ := password.Generate(2, 2, 0, false, false)
	a2b, _ := password.Generate(2, 0, 0, true, false)
	d2b, _ := password.Generate(2, 2, 0, false, false)

	return a2 + d2 + a2b + d2b
}

// OutRegion is a Region for outputting JSON hierarchy
type OutRegion struct {
	Name    string       `json:"name"`
	Xid     string       `json:"xid"`
	Schools []*OutSchool `json:"schools"`
}

func (r *OutRegion) getSchool(school string, town string) (*OutSchool, error) {
	ErrNotFound := errors.New("not Found")
	for _, val := range r.Schools {
		if val.Name == school && val.Town == town {
			return val, nil
		}
	}
	return nil, ErrNotFound
}
func (r *OutRegion) makeXIDS() {
	r.Xid = makeXID()

	for _, school := range r.Schools {
		school.Xid = makeXID()
	}

}

// ByName provides the interface for doing built-in sorting of schools within a region
type ByName []*OutSchool

// Len is the length of the array
func (a ByName) Len() int { return len(a) }

// Less is whether I should sort before j
func (a ByName) Less(i, j int) bool { return a[i].Name+a[i].Town < a[j].Name+a[j].Town }

// Swap swaps i and j
func (a ByName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

//OutSchool is a School for outputting JSON hierarchy
type OutSchool struct {
	Name string `json:"name"`
	Xid  string `json:"xid"`
	Town string `json:"town,omitempty"`
}

// SchoolCsv is the flattened (with region) version of a school for csv output and sorting
type SchoolCsv struct {
	Name   string
	Region string
	Xid    string
	Town   string
}

func (s *SchoolCsv) asCsv() string {
	if s.Town == "" {
		return fmt.Sprintf("\"%v\",\"%v\",\"%v\"", s.Xid, s.Name, s.Region)
	}
	return fmt.Sprintf("\"%v\",\"%v (%v)\",\"%v\"", s.Xid, s.Name, s.Town, s.Region)
}

// BySchool provides the interface for doing built-in sorting of schools
type BySchool []*SchoolCsv

// Len is the length of the array
func (a BySchool) Len() int { return len(a) }

// Less is whether I should sort before j
func (a BySchool) Less(i, j int) bool { return a[i].Name+a[i].Town < a[j].Name+a[j].Town }

// Swap swaps i and j
func (a BySchool) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func main() {
	// urlbase := "http://localhost:8080"
	// rulesLocation := "http://www.cteresource.org/fbla/hs_competitive_events/index.html"
	schoolFilesLoc := "FBLA_2019/exams/regional/test_papers/*.json"
	// examTime := 50
	// contestInfo := "exams/regional/contest_info/contest_info.json"
	schoolFiles, _ := filepath.Glob(schoolFilesLoc)
	// loadExams := true

	var state State

	state.Id = "VA"
	state.Name = "Virginia"

	var csvFile strings.Builder

	eventout := OutEvent{Name: "2019 Regionals"}

	var town string

	var schoolsCSV []*SchoolCsv

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
		outpt := "region: " + testPapers.Region + " school: " + testPapers.School

		region, err := eventout.getRegion(testPapers.Region)

		if err != nil {
			newRegion := OutRegion{Name: testPapers.Region}
			eventout.Regions = append(eventout.Regions, &newRegion)
			region = &newRegion
		}
		// eventout.Regions = append(eventout.Regions, region)

		school, err := region.getSchool(testPapers.School, town)

		if err != nil {

			newSchool := OutSchool{Name: testPapers.School}
			if town != "" {
				newSchool.Town = town
			}
			school = &newSchool
		}
		region.Schools = append(region.Schools, school)

		if len(town) > 0 {
			outpt = outpt + " town: " + town
		}
		fmt.Println(outpt)

	}
	eventout.makeXIDS()
	for _, region := range eventout.Regions {
		sort.Sort(ByName(region.Schools))
	}
	for _, region := range eventout.Regions {
		r2 := Region{Id: region.Xid, Name: region.Name}
		state.Regions = append(state.Regions, &r2)
		for _, school := range region.Schools {
			s2 := School{Name: school.Name, Town: school.Town, Id: school.Xid}
			r2.Schools = append(r2.Schools, &s2)
			schoolsCSV = append(schoolsCSV, &SchoolCsv{Xid: school.Xid, Name: school.Name, Region: region.Name, Town: school.Town})

		}
	}
	// CSV
	sort.Sort(BySchool(schoolsCSV))
	for _, school := range schoolsCSV {
		fmt.Fprintf(&csvFile, "%s\n", school.asCsv())
	}
	ioutil.WriteFile("schools.csv", []byte(csvFile.String()), 0644)

	// JSON
	file, _ := json.MarshalIndent(eventout, "", " ")
	ioutil.WriteFile("schools.json", file, 0644)

	// Protocol Buffers
	tfile, _ := proto.Marshal(&state)

	ioutil.WriteFile("schools.pb", tfile, 0644)

	fmt.Printf("MAC: %16.16X\n", macUint64())

	var newstate proto.Message
	err := proto.Unmarshal(tfile, newstate)
	if err != nil {
		fmt.Println(err)
	}
}
