package exception

type ConflictError struct {
	Error string
}

func NewConflictError(error string) ConflictError {
	return ConflictError{Error: error}
}
