package net

import "io"

type WriteCloser interface {
	io.WriteCloser
}
