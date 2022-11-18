package models

// ClassworkEntry struct to represent all classswork for a single class
// during a given six weeks.
type ClassworkEntry struct {
	Position    int          `json:"position"`    //The position of the class, used for ordering
	Class       Class        `json:"class"`       //Class information about the entry
	Average     int          `json:"average"`     //The average grade for that class
	Assignments []Assignment `json:"assignments"` //All the assignments currently entered for the class
}

// Classwork represents all classwork for a
// specific six weeks, stored in an array.
type Classwork struct {
	SixWeeks int              `json:"sixWeeks"` //The six weeks the classwork is for
	Entries  []ClassworkEntry `json:"entries"`  //An array of ClassworkEntry structs containing classwork for each class
}
