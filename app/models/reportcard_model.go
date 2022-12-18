package models

// RequestCardRequestBody represents the
// request body that is to be passed in
// with the POST request to this endpoint.
type ReportCardRequestBody struct {
	BaseRequestBody
}

// SixWeeksOther contains fields for
// only every six weeks, and is used
// to contain miscellaneous data such as
// comments or absences.
type SixWeeksOther struct {
	First  string `json:"first"`
	Second string `json:"second"`
	Third  string `json:"third"`
	Fourth string `json:"fourth"`
	Fifth  string `json:"fifth"`
	Sixth  string `json:"sixth"`
}

// SixWeeksWithExam contains fields for every
// six weeks plus exams/semesters, used to
// contain grades/averages.
type SixWeeksGrades struct {
	SixWeeksOther
	Exam1 string `json:"exam1"`
	Sem1  string `json:"sem1"`
	Exam2 string `json:"exam2"`
	Sem2  string `json:"sem2"`
}

// Absences represents the struct
// specifically made for the absences
// data present in a report card.
type Absences struct {
	ExcusedAbsence   string `json:"excusedAbsence"`   // The amount of excused absences for the class
	UnexcusedAbsence string `json:"unexcusedAbsence"` // The amount of unexcused absences for the class
	ExcusedTardy     string `json:"excusedTardy"`     // The amount of excused tardies for the class
	UnexcusedTardy   string `json:"unexcusedTardy"`   // The amount of unexcused tardies for the class
}

// ReportCardEntry represents a singular
// entry in the report card for one class.
type ReportCardEntry struct {
	Class           Class          `json:"class"`           // Information about the class for the entry
	AttemptedCredit string         `json:"attemptedCredit"` // The amount of credit attempted
	EarnedCredit    string         `json:"earnedCredit"`    // The amount of credit earned
	Averages        SixWeeksGrades `json:"averages"`        // Data about grades
	Comments        SixWeeksOther  `json:"comments"`        // Data about comments
	Conduct         SixWeeksOther  `json:"conduct"`         // Data about conduct
	Absences        Absences       `json:"absences"`        // Data about absences
}

// ReportCard holds the array Entries with
// the report card data for each class.
type ReportCard struct {
	Entries []ReportCardEntry `json:"entries"` // All the report card entries
}

// ReportCardResponse represents a JSON response
// to the Report Card POST request.
type ReportCardResponse struct {
	HTTPError               // Error, if one is attached to the response
	ReportCard []ReportCard `json:"reportCard"` // The resulting report card
}
