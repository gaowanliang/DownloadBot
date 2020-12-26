package http

import "errors"

// error types

var (
	ErrFileTooLarge = errors.New("file is too large for simple upload")
)

type innerError struct {
	Code       string      `json:"code"`
	Message    string      `json:"message"`
	InnerError *innerError `json:"innererror"`
}

// The Error type defines the basic structure of errors that are returned from
// the OneDrive API.
// See: http://onedrive.github.io/misc/errors.htm
type Error struct {
	innerError `json:"error"`
}

func (e Error) Error() string {
	return e.Message
}
