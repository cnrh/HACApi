package models

// BaseRequestBody describes a struct with the base properties needed
// for most POST request bodies.
type BaseRequestBody struct {
	// The username to log in with
	Username string `json:"username" validate:"required,min=1" example:"j1732901"`
	// The password to log in with
	Password string `json:"password" validate:"required,min=1" example:"j382704"`
	// The base URL for the PowerSchool HAC service
	Base string `json:"base" validate:"required,min=1" example:"https://homeaccess.katyisd.org"`
}
