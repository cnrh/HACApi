package models

// Assignment struct to represent a singular assignment entry in HAC.
type Assignment struct {
	DueDate      string `json:"dueDate"`      //The date the assignment is due
	AssignedDate string `json:"assignedDate"` //The date the assignment was assigned
	Name         string `json:"name"`         //The name of the assignment
	Category     string `json:"category"`     //The category of the assignment (major, minor, other, etc)
	Grade        string `json:"grade"`        //What grade the user got on the assignment
	TotalPoints  string `json:"totalPoints"`  //The total points that could be earned on the assignment
	Dropped      bool   `json:"dropped"`      //Whether the assignment was dropped or not
}
