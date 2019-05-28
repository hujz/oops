package errors

const Err_Find_Dependency string = "Err_Find_Dependency"

func NewError(code, message string) error {
	return &errorString{code, message}
}

type errorString struct {
	c, m string
}

func (e *errorString) Error() string {
	return e.c + ":" + e.m
}
