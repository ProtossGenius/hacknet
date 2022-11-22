package gkritf

import "net"

// StringAble can to string and from string.
type StringAble interface {
	String() string
	From(str string)
	Inited() bool
}

// GeekerNodeInfo a not service's info.
type GeekerNodeInfo struct {
	StringAble
	IP   string
	Port int
}

// BuildUDPAddr build udp addr.
func (gi GeekerNodeInfo) BuildUDPAddr() *net.UDPAddr {
	return &net.UDPAddr{
		IP:   net.ParseIP(gi.IP),
		Port: gi.Port,
	}
}

// IGeekerNetUDPServer geeker server.
type IGeekerNetUDPServer interface {
	// Listen listion a local udp port.
	Listen(port int) error
	// Search search a seesionId's info.
	Search(sessionID string) GeekerNodeInfo
	// Close close.
	Close()
}

// ISessionManager session manager.
type ISessionManager interface {
	// Register register a session.
	Register(sessionID string, info GeekerNodeInfo) error
	// Search search a sessionId.
	Search(sessionID string) GeekerNodeInfo
}
