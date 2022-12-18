package models

// TranscriptRequestBody is the body to be passed with
// a POST request to the transcript endpoint
type TranscriptRequestBody struct {
	BaseRequestBody
}

// TranscriptGroupEntry represents a singular entry
// in a transcript group.
type TranscriptGroupEntry struct {
	Class   Class  `json:"class"`   // The class related to the entry
	Average string `json:"average"` // The average grade for the class
	Credit  string `json:"credit"`  // The credit earned for the class
}

// TranscriptGroup represents a group of entries,
// usually grouped by term.
type TranscriptGroup struct {
	Year        string                 `json:"year"`        // The school year the group represents data for
	Semester    string                 `json:"semester"`    // The semester the group represents data for
	GradeLevel  string                 `json:"gradeLevel"`  // The grade level the group represents data for
	Building    string                 `json:"building"`    // The building in which the classes in the group weree taken
	Entries     []TranscriptGroupEntry `json:"entries"`     // The individual class entries for this group
	TotalCredit string                 `json:"totalCredit"` // The total credit earned for all the entries in this group
}

// TranscriptGPA describes information about weighted/unweighted GPA.
// for the user.
type TranscriptGPA struct {
	Type     string `json:"type"`     // The type of GPA (unweighted/weighted)
	GPA      string `json:"gpa"`      // The GPA
	Rank     string `json:"rank"`     // The class rank
	Quartile string `json:"quartile"` // The class quartile
}

// Transcript represents all the transcript groups and both types of GPA.
type Transcript struct {
	Entries    []TranscriptGroup `json:"entries"`    // All the transcript groups
	Weighted   TranscriptGPA     `json:"weighted"`   // The weighted GPA
	Unweighted TranscriptGPA     `json:"unweighted"` // The unweighted GPA
}

// TranscriptResponse represents a JSON response
// to the Transcript POST request.
type TranscriptResponse struct {
	HTTPError               // Error, if one is attached to the response
	Transcript []Transcript `json:"transcript"` // The resulting transcript
}
