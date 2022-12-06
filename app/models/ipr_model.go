package models

// IprRequestBody represents the body that is to
// be passed with a POST request to the IPR
// endpoint.
type IprRequestBody struct {
	BaseRequestBody
	// The date of the IPR to return
	Date string `json:"date" example:"09/06/2022"`
}

// IPREntry represents an individual class's progress
// report in the overall progress report.
type IPREntry struct {
	Class Class  `json:"class"` // Information about the class related to the IPREntry
	Grade string `json:"grade"` // The average at the moment the progress report was submitted
}

// IPR represents the overall progress report. It contains
// an array containing all the IPR entries.
type IPR struct {
	Date    string     `json:"date"`    // The date the IPR was submitted
	Entries []IPREntry `json:"entries"` // An array representing all the IPR entries
}

// IPRResponse represents a JSON response
// to the IPR POST request.
type IPRResponse struct {
	HTTPError       // Error, if one is attached to the response
	IPR       []IPR `json:"ipr"` // The resulting IPR(s)
}
