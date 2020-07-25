package main

type ExamInfo struct {
	Page2Pdf string `json:"page2pdf"`
	Manual   string `json:"manual"`
	// PngOffset   string   `json:"pngOffset"`
	GradeLevels []string `json:"grade_levels"`
	Level       string   `json:"level"`
	Exams       []Exams  `json:"exams"`
}
type Exams struct {
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

type ExamV2 struct {
	Items []ExamItem `json:"items"`
	// Time  int    `json:"time"`
	ID    string `json:"id"`
	Title string `json:"title"`
}

type Choices struct {
	ID       string `json:"id"`
	Response string `json:"response"`
}

type ExamItem struct {
	ID            string    `json:"id"`
	CorrectAnswer string    `json:"correct_answer"`
	Question      string    `json:"question"`
	Rubric        string    `json:"rubric"`
	QuestionType  string    `json:"question_type"`
	Choices       []Choices `json:"choices"`
}
type sortedItems []ExamItem

func (s sortedItems) Len() int {
	return len(s)
}

func (s sortedItems) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s sortedItems) Less(i, j int) bool {
	return s[i].ID < s[j].ID
}
