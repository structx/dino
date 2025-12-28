package net

import (
	"io"
)

type ReadCloser interface {
	io.ReadCloser
}
