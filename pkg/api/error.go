package api

// Error returns a descriptive string.
func (e *Issue) Error() string {
	return e.Message
}

// Append adds an Issue to Error.
func (e *Error) Append(err *Issue) *Error {
	e.Issues = append(e.Issues, err)
	return e
}

// Error returns a descriptive string.
func (e *Error) Error() string {
	return e.Message
}
