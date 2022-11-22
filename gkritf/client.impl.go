package gkritf

import (
	"fmt"
	"net"
)

// GeekerNetUDPClient geeker net udp client.
type GeekerNetUDPClient struct {
	conn      *net.UDPConn
	sessionID string
	cmdMap    map[string]FCommand
}

// NewGeekerNetUDPClient create new client.
func NewGeekerNetUDPClient(sessionID string) IGeekerNetUDPClient {
	return (&GeekerNetUDPClient{sessionID: sessionID, cmdMap: make(map[string]FCommand)})
}

// Connect connect to a server.
func (g *GeekerNetUDPClient) Connect(localPort int, remoteIP string, remotePort int) (err error) {
	localAddr := &net.UDPAddr{IP: net.IPv4zero, Port: localPort}
	remoteAddr := &net.UDPAddr{IP: net.ParseIP(remoteIP), Port: remotePort}
	if g.conn, err = net.DialUDP("udp", localAddr, remoteAddr); err != nil {
		fmt.Println(localAddr.String(), " connect ", remoteAddr.String(), " fail, error is ", err)

		return err
	}

	go g.startListen()

	return g.Send("Register#" + g.sessionID)
}

func (g *GeekerNetUDPClient) startListen() {
	data := make([]byte, 1024)
	for {
		n, remoteAddr, err := g.conn.ReadFromUDP(data)
		if err != nil {
			fmt.Printf("error during read: %s", err)
		}
		fmt.Println("read something from ", remoteAddr.String(),
			", msg := ", string(data[:n]),
			", error := ", err)
	}
}

// Send send message to GeekerNetUDPServer.
func (g *GeekerNetUDPClient) Send(msg string) (err error) {
	fmt.Println("send message :", msg)
	if _, err = g.conn.Write([]byte(msg)); err != nil {
		fmt.Println(g.sessionID, " send msg ", msg, ",get error :", err)
	}

	return
}

// RegisterMsgHandler register what will do when get a command.
func (g *GeekerNetUDPClient) RegisterMsgHandler(name string, fc FCommand) {
	g.cmdMap[name] = fc
}
