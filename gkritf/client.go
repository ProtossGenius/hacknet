package gkritf

// IGeekerNetUDPClient geeker net udp client.
type IGeekerNetUDPClient interface {
	// Connect connect to a server.
	Connect(localPort int, remoteIP string, remotePort int) error
	// Send send message to GeekerNetUDPServer.
	Send(msg string) error
	// RegisterMsgHandler register what will do when get a command.
	RegisterMsgHandler(name string, fc FCommand)
}
