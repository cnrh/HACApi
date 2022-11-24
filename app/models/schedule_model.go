package models

// ScheduleEntry represents a singular entry in the schedule.
type ScheduleEntry struct {
	Class          Class    // Information about the Class related to the Schedule
	Days           []string //The days the class is active for
	Building       string   //The building the class is in
	Active         bool     //Whether the class is active or not
	MarkingPeriods []string //The marking periods the class is active for
}

// Schedule contains an array with all the
// schedule entries for the user
type Schedule struct {
	Entries []ScheduleEntry //An array containing all the schedule entries
}
