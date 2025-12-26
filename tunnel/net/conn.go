package net

// Conn
type Conn interface {
	Read() (*DataFrame, error)
	Write(*DataFrame) (int, error)

	Close(string) error
}

// ClientConn
type ClientConn interface {
	Read() (*ClientFrame, error)
	Write(*ClientFrame) (int, error)

	Close(string) error
}
