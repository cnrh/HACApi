package models

// ClassworkRequestBody represents the body that is to be passed with
// the POST request to the classwork endpoint.
type ClassworkRequestBody struct {
	BaseRequestBody
	// The marking period to pull data from
	MarkingPeriods []int `json:"markingPeriods" validate:"max=6,dive,min=1,max=6" example:"1,2"`
}

// ClassworkEntry represents all classswork for a single class
// during a given six weeks.
type ClassworkEntry struct {
	Position    int          `json:"position"`    // The position of the class, used for ordering
	Class       Class        `json:"class"`       // Class information about the entry
	Average     string       `json:"average"`     // The average grade for that class
	Assignments []Assignment `json:"assignments"` // All the assignments currently entered for the class
}

// Classwork represents all classwork
// for a specific six weeks, stored in an array.
type Classwork struct {
	MarkingPeriod int              `json:"sixWeeks"` // The marking period the classwork is for
	Entries       []ClassworkEntry `json:"entries"`  // An array of ClassworkEntry structs containing classwork for each class
}

// ClassworkResponse represents a JSON response
// to the Classwork POST request.
type ClassworkResponse struct {
	HTTPError             // Error, if one is attached to the response
	Classwork []Classwork `json:"classwork"` // The resulting classwork
}
