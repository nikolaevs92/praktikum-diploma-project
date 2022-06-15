package statuserror

type StatusError struct {
	Status int
	Err    string
}

func (st StatusError) Error() string {
	return st.Err
}
func NewStatusError(err string, status int) StatusError {
	return StatusError{Status: status, Err: err}
}
