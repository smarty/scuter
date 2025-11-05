package scuter

// Error represents some kind of problem, most likely with the calling HTTP request.
type Error struct {
	// Fields indicates the exact location(s) of the errors including the part of
	// the HTTP request itself this is invalid. Valid field prefixes include
	// "path", "query", "header", "form", and "body".
	Fields []string `json:"fields,omitempty"`

	// ID represents the unique, numeric contractual identifier that can be used to
	// associate this error with a particular front-end error message, if any.
	ID int `json:"id,omitempty"`

	// Name represents the unique string-based, contractual value that can be used to
	// associate this error with a particular front-end error message, if any.
	Name string `json:"name,omitempty"`

	// Message represents a friendly, user-facing message to indicate why there was a
	// problem with the input.
	Message string `json:"message,omitempty"`
}

func (this Error) Error() string { return this.Message }

// Errors represents a set of problems.
type Errors struct {
	Errors []Error `json:"errors,omitempty"`
}

func NewErrors(values ...Error) *Errors {
	return &Errors{Errors: values}
}

func (this *Errors) Append(errs ...Error) {
	this.Errors = append(this.Errors, errs...)
}
