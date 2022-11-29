package models

// TranscriptGroupEntry represents a singular entry
// in a transcript group.
type TranscriptGroupEntry struct {
	Class   Class  // The class related to the entry
	Average string // The average grade for the class
	Credit  string // The credit earned for the class
}

// TranscriptGroup represents a group of entries,
// usually grouped by term.
type TranscriptGroup struct {
	Year        string                 // The school year the group represents data for
	Semester    string                 // The semester the group represents data for
	GradeLevel  string                 // The grade level the group represents data for
	Building    string                 // The building in which the classes in the group weree taken
	Entries     []TranscriptGroupEntry // The individual class entries for this group
	TotalCredit string                 // The total credit earned for all the entries in this group
}

// TranscriptGPA describes information about weighted/unweighted GPA.
// for the user.
type TranscriptGPA struct {
	Type     string // The type of GPA (unweighted/weighted)
	GPA      string // The GPA
	Rank     string // The class rank
	Quartile string // The class quartile
}

// Transcript represents all the transcript groups and both types of GPA.
type Transcript struct {
	Groups     []TranscriptGroup // All the transcript groups
	Weighted   TranscriptGPA     // The weighted GPA
	Unweighted TranscriptGPA     // The unweighted GPA
}
