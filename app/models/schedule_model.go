package models

// ScheduleRequestBody represents the
// request body to be passed in with a
// POST request to the endpoint.
type ScheduleRequestBody struct {
	BaseRequestBody
}

// ScheduleEntry represents a singular entry in the schedule.
type ScheduleEntry struct {
	Class          Class    `json:"class"`          // Information about the Class related to the Schedule
	Days           []string `json:"days"`           // The days the class is active for
	Building       string   `json:"building"`       // The building the class is in
	Active         bool     `json:"active"`         // Whether the class is active or not
	MarkingPeriods []string `json:"markingPeriods"` // The marking periods the class is active for
}

// Schedule contains an array with all the
// schedule entries for the user
type Schedule struct {
	Entries []ScheduleEntry `json:"entries"` // An array containing all the schedule entries
}

// ScheduleResponse represents a JSON response
// to the Schedule POST request.
type ScheduleResponse struct {
	HTTPError          // Error, if one is attached to the response
	Schedule  Schedule `json:"schedule"` // The resulting schedule
}
