package models

// HTTPError represents a HTTP error.
type HTTPError struct {
	Error   bool   `json:"err"` // If there was an error
	Message string `json:"msg"` // The associated message
}
