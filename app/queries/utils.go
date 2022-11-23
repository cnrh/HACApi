package queries

// PipelineResponse represents a intermediary
// response between pipelines, used extensively
// in the queries package.
type PipelineResponse[T any] struct {
	Value T
	Err   error
}

// partialFormData struct representing partial formdata
// for a POST request.
type partialFormData struct {
	ViewState       string //viewstate formdata entry
	ViewStateGen    string //viewstategen formdata entry
	EventValidation string //eventvalidation formdata entry
	Url             string //url for the request
	Base            string //base url for the request
}
