package utils

// HTTPError struct representing a HTTP error
type HTTPError struct {
	Code    int    `json:"code"`    //The HTTP code
	Message string `json:"message"` //The associated message
}
