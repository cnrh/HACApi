package models

// Class represents a class in HAC.
type Class struct {
	Name    string `json:"name"`    // The name of the class
	Course  string `json:"course"`  // The course ID of the class
	Period  string `json:"period"`  // What period the class is for the student, relative to the current schedule
	Teacher string `json:"teacher"` // The name of the teacher of the class
	Room    string `json:"room"`    // The room number of the class
}
