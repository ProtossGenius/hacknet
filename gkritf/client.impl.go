package gkritf

import (
	"fmt"
	"net"
	"strings"
)

// GeekerNetUDPClient geeker net udp client.
type GeekerNetUDPClient struct {
	conn       *net.UDPConn
	sessionID  string
	cmdMap     map[string]FCommand
	remoteAddr *net.UDPAddr
}

// GetConn only get conn
func (g *GeekerNetUDPClient) GetConn() *net.UDPConn {
	return g.conn
}

// MoveConn get net.UDPConn and remove from client(then stop the client).
func (g *GeekerNetUDPClient) MoveConn() *net.UDPConn {
	conn := g.conn
	g.conn = nil
	return conn
}

// NewGeekerNetUDPClient create new client.
func NewGeekerNetUDPClient(sessionID string) IGeekerNetUDPClient {
	return (&GeekerNetUDPClient{sessionID: sessionID, cmdMap: make(map[string]FCommand)})
}

// Connect connect to a server.
func (g *GeekerNetUDPClient) Connect(localPort int, remoteIP string, remotePort int) (err error) {
	localAddr := &net.UDPAddr{IP: net.IPv4zero, Port: localPort}
	remoteAddr := &net.UDPAddr{IP: net.ParseIP(remoteIP), Port: remotePort}
	g.remoteAddr = remoteAddr
	if g.conn, err = net.ListenUDP("udp", localAddr); err != nil {
		fmt.Println(localAddr.String(), " connect ", remoteAddr.String(), " fail, error is ", err)

		return err
	}

	go g.startListen()
	fmt.Println("client ", g.sessionID, " startup, listen port ", localPort, ".")
	return g.Send(RegisterC{SessionID: g.sessionID}.Message())
}

func (g *GeekerNetUDPClient) startListen() {
	data := make([]byte, 1024)

	unhandle := func(addr *net.UDPAddr, msg GeekerMsg) {
		fmt.Println("unhandle message from ", addr.String(),
			", msg := ", msg)
	}

	for conn := g.conn; conn != nil; conn = g.conn {
		n, remoteAddr, err := g.conn.ReadFromUDP(data)
		if err != nil {
			fmt.Printf("error during read: %s", err)
		}
		msg := string(data[:n])

		list := strings.SplitN(msg, "#", 2)
		if len(list) < 2 {
			unhandle(remoteAddr, GeekerMsg(msg))
			continue
		}

		if ff, ok := g.cmdMap[list[0]]; ok {
			ff(remoteAddr, GeekerMsg(list[1]))
		} else {
			unhandle(remoteAddr, GeekerMsg(msg))
		}

	}
}

// Send send message to GeekerNetUDPServer.
func (g *GeekerNetUDPClient) Send(msg GeekerMsg) (err error) {
	fmt.Println("send message :", msg)
	if _, err = g.conn.WriteToUDP([]byte(msg), g.remoteAddr); err != nil {
		fmt.Println(g.sessionID, " send msg ", msg, ",get error :", err)
	}

	return
}

// RegisterMsgHandler register what will do when get a command.
func (g *GeekerNetUDPClient) RegisterMsgHandler(name string, fc FCommand) {
	g.cmdMap[name] = fc
}
