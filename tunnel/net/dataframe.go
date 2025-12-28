package net

// DataFrame tunnel data frame
type DataFrame struct {
	SessionID      string
	NewConn        *NewConn
	CloseConn      *CloseConn
	Payload        []byte
	RouteUpdate    *RouteUpdate
	IsControlFrame bool
}

// NewConn
type NewConn struct {
	Hostname string
}

// CloseConn
type CloseConn struct {
	Status int
}

type RouteUpdate struct {
	Hostname     string
	DestProtocol string
	DestIP       string
	DestPort     uint32
	IsDelete     bool
}
