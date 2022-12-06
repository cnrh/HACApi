package models

// LoginRequestBody represents the body that is to be
// passed along with the POST request to the login
// endpoint.
type LoginRequestBody struct {
	BaseRequestBody
}

// Login represents the response to the POST request to
// the login endpoint.
type Login struct {
	Username string `json:"username"` // The username used to sign in with
	Password string `json:"password"` // The password used to sign in with
	Base     string `json:"base"`     // The base URL signed in to
}

// LoginResponse represents a JSON response
// to the Login POST request.
type LoginResponse struct {
	HTTPError       // Error, if one is attached to the response
	Login     Login `json:"login"` // Data about the login
}
