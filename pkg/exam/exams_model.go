package exam

// RegistrationMeta is an umbrella struct for metadata about delivering a set of exams
type RegistrationMeta struct {
	GradeLevels []string `json:"grade_levels"`
	Level       string   `json:"level"`
	Exams       []Meta   `json:"exams"`
}

// Meta is a metadata about registering an exam
type Meta struct {
	// PdfPage          string `json:"pdfPage"`
	MaxTeams         int    `json:"maxTeams"`
	Title            string `json:"title"`
	CollaborationMax int    `json:"collaborationMax"`
	MinGradeLevel    string `json:"minGradeLevel"`
	Registration     string `json:"registration"`
	ID               string `json:"id"`
	MaxGradeLevel    string `json:"maxGradeLevel"`
	Filename         string `json:"filename"`
	Message          string `json:"message,omitempty"`
	PapersPerTeam    int    `json:"papersPerTeam,omitempty"`
	NeedGradeLevel   bool   `json:"needGradeLevel,omitempty"`
}

// Examination holds the info for an exam
type Examination struct {
	Items []Item `json:"items"`
	ID    string `json:"id"`
	Title string `json:"title"`
}

// Choice is id:response for a choice in an exam
type Choice struct {
	ID       string `json:"id"`
	Response string `json:"response"`
}

// Item is an exam item. A question on the exam.
type Item struct {
	ID            string   `json:"id"`
	CorrectAnswer string   `json:"correct_answer"`
	Question      string   `json:"question"`
	Rubric        string   `json:"rubric"`
	QuestionType  string   `json:"question_type"`
	Choices       []Choice `json:"choices"`
}

// ItemsByID sorts a list of Items by ID
type ItemsByID []Item

func (s ItemsByID) Len() int {
	return len(s)
}

func (s ItemsByID) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ItemsByID) Less(i, j int) bool {
	return s[i].ID < s[j].ID
}
