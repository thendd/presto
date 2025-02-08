package errors

import "errors"

var UnknwonError = errors.New("There was an unexpected error while executing this command")

// This function was created to avoid confusion between this local `errors` package
// and the one provided by Go
func New(text string) error {
	return errors.New(text)
}
