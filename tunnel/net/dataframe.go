package net

// DataFrame tunnel data frame
type DataFrame struct {
	SessionID      string
	NewConn        *NewConn
	Payload        []byte
	RouteUpdate    *RouteUpdate
	IsControlFrame bool
}

// ClientFrame
type ClientFrame struct {
	SessionID      string
	NewConn        *NewConn
	Payload        []byte
	RouteUpdate    *RouteUpdate
	IsControlFrame bool
}

type NewConn struct {
	Hostname string
}

type RouteUpdate struct {
	Hostname     string
	DestProtocol string
	DestIP       string
	DestPort     uint32
	IsDelete     bool
}
