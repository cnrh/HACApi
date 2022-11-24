package models

// SixWeeksWithExam contains fields for every
// six weeks plus exams/semesters, used to
// contain grades/averages.
type SixWeeksGrades struct {
	First  string
	Second string
	Third  string
	Exam1  string
	Sem1   string
	Fourth string
	Fifth  string
	Sixth  string
	Exam2  string
	Sem2   string
}

// SixWeeksOther contains fields for
// only every six weeks, and is used
// to contain miscellaneous data such as
// comments or absences.
type SixWeeksOther struct {
	First  string
	Second string
	Third  string
	Fourth string
	Fifth  string
	Sixth  string
}

// Absences represents the struct
// specifically made for the absences
// data present in a report card.
type Absences struct {
	ExcusedAbsence   string
	UnexcusedAbsence string
	ExcusedTardy     string
	UnexcusedTardy   string
}

// ReportCardEntry represents a singular
// entry in the report card for one class.
type ReportCardEntry struct {
	Class           Class
	AttemptedCredit string
	EarnedCredit    string
	Averages        SixWeeksGrades
	Comments        SixWeeksOther
	Conduct         SixWeeksOther
	Absences        Absences
}

// ReportCard holds the array Entries with
// the report card data for each class.
type ReportCard struct {
	Entries []ReportCardEntry
}
