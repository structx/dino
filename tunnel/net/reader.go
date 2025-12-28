package net

import (
	"io"
	"net"
	"sync"
)

type Reader interface {
	io.Reader
}

type tunnelReader struct {
	close    sync.Once
	isclosed bool
	inbound  chan []byte
}

// NewReader
func NewReader(inboundCh chan []byte) Reader {
	return &tunnelReader{
		close:    sync.Once{},
		isclosed: false,
		inbound:  inboundCh,
	}
}

// Close implements ReadCloser.
func (t *tunnelReader) Close() error {
	t.close.Do(func() {
		t.isclosed = true
	})
	return nil
}

// Read implements ReadCloser.
func (t *tunnelReader) Read(p []byte) (n int, err error) {
	if t.isclosed {
		return 0, net.ErrClosed
	}

	data, ok := <-t.inbound
	if !ok {
		return 0, net.ErrClosed
	}

	return copy(p, data), nil
}
