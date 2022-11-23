package models

// ClassworkResponse represents a JSON response
// to the Classwork POST request.
type ClassworkResponse struct {
	HTTPError             //Error, if one is attached to the response
	Classwork []Classwork `json:"classwork"` //The resulting classwork
}

// LoginResponse represents a JSON response
// to the Login POST request.
type LoginResponse struct {
	HTTPError //Error, if one is attached to the response
}

// IPRResponse represents a JSON response
// to the IPR POST request
type IPRResponse struct {
	HTTPError       //Error, if one is attached to the response
	IPR       []IPR `json:"ipr"` //The resulting IPR(s)
}
