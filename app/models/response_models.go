package models

// ClassworkResponse represents a JSON response
// to the Classwork POST request.
type ClassworkResponse struct {
	HTTPError           //Error, if one is attached to the response
	Classwork Classwork `json:"classwork"` //The resulting classwork
}

// LoginResponse represents a JSON response
// to the Login POST request.
type LoginResponse struct {
	HTTPError //Error, if one is attached to the response
}
