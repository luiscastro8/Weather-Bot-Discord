package testutil

import (
	"errors"
	"io"
)

type errReader struct{}

func (e errReader) Read(p []byte) (int, error) {
	return 0, errors.New("intentional error")
}

func NewErrorReader() io.Reader {
	return errReader{}
}
