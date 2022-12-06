package models

// IprAllRequestBody represents the body that is to
// be passed with a POST request to the /ipr/all
// endpoint.
type IprAllRequestBody struct {
	BaseRequestBody
	// Whether to return only dates or all the IPRs
	DatesOnly bool `json:"datesOnly" example:"true" default:"false"`
}
