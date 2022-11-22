package gkritf

import "net"

// IGeekerNetUDPClient geeker net udp client.
type IGeekerNetUDPClient interface {
	// Connect connect to a server.
	Connect(localPort int, remoteIP string, remotePort int) error
	// Send send message to GeekerNetUDPServer.
	Send(msg GeekerMsg) error
	// RegisterMsgHandler register what will do when get a command.
	RegisterMsgHandler(name string, fc FCommand)
	// MoveConn get net.UDPConn and remove from client(then stop the client).
	MoveConn() *net.UDPConn
	// GetConn only get conn
	GetConn() *net.UDPConn
}
