package net

// Conn
type Conn interface {
	Read() (*DataFrame, error)
	Write(*DataFrame) (int, error)

	Close(string) error
}
