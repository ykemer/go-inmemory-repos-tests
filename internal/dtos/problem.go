package dtos

// ProblemDetail is the RFC 7807 error response shape returned by all error
// endpoints. It exists here so swag can generate accurate schema references.
type ProblemDetail struct {
	Type          string            `json:"type"`
	Title         string            `json:"title"`
	Status        int               `json:"status"`
	Detail        string            `json:"detail,omitempty"`
	InvalidParams map[string]string `json:"invalid_params,omitempty"`
}
