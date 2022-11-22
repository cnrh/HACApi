package queries

// PipelineResponse represents a intermediary
// response between pipelines, used extensively
// in the queries package.
type PipelineResponse[T any] struct {
	Value T
	Err   error
}
