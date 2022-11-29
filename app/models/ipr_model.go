package models

// IPREntry represents an individual class's progress
// report in the overall progress report.
type IPREntry struct {
	Class Class  // Information about the class related to the IPREntry
	Grade string // The average at the moment the progress report was submitted
}

// IPR represents the overall progress report. It contains
// an array containing all the IPR entries.
type IPR struct {
	Date    string     // The date the IPR was submitted
	Entries []IPREntry // An array representing all the IPR entries
}
