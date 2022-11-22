package gkritf

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

// Inited is info inited.
func (gi GeekerNodeInfo) Inited() bool {
	return gi.IP != "" && gi.Port != 0
}
func (gi GeekerNodeInfo) String() string {
	return fmt.Sprintf("%s#%d", gi.IP, gi.Port)
}

// From init from string.
func (gi GeekerNodeInfo) From(str string) {
	arr := strings.Split(str, "#")
	gi.IP = arr[0]
	gi.Port, _ = strconv.Atoi(arr[1])
}

// NewGeekerNetUDPServer new geeker net udp server.
func NewGeekerNetUDPServer() IGeekerNetUDPServer {
	return (&GeekerNetUDPServer{
		sessionMgr: &sessionManager{
			sessionMap: make(map[string]GeekerNodeInfo, 0),
		},
		running: false,
	}).init()
}

// FCommand what will command do.
type FCommand func(addr *net.UDPAddr, params string) error

// GeekerNetUDPServer udp server.
type GeekerNetUDPServer struct {
	sessionMgr ISessionManager
	cmdMap     map[string]FCommand
	running    bool
	listener   *net.UDPConn
}

// Close close.
func (g *GeekerNetUDPServer) Close() {
	g.running = false
}

// Listen listion a local udp port.
func (g *GeekerNetUDPServer) init() IGeekerNetUDPServer {

	sendmsg := func(addr *net.UDPAddr, msg string) (err error) {
		if _, err = g.listener.WriteToUDP([]byte(msg), addr); err != nil {
			log.Println("when send response msg : ", msg, ", error is ", err)
		}

		return err
	}

	g.cmdMap = map[string]FCommand{
		"Search": func(addr *net.UDPAddr, sessionID string) error {
			return sendmsg(addr, g.Search(sessionID).String())
		},
		"Register": func(addr *net.UDPAddr, sessionID string) error {
			g.sessionMgr.Register(sessionID, GeekerNodeInfo{IP: addr.IP.String(), Port: addr.Port})

			return sendmsg(addr, "Register:done")
		},
		"Notice": func(addr *net.UDPAddr, noticeInfo string) error {
			nodeInfo := GeekerNodeInfo{IP: addr.IP.String(), Port: addr.Port}
			sendmsg(nodeInfo.BuildUDPAddr(), "Notice#"+nodeInfo.String()+"#"+noticeInfo)
			return sendmsg(addr, "Notice:done")
		},
	}

	return g
}

// Listen listion a local udp port.
func (g *GeekerNetUDPServer) Listen(port int) error {
	if g.running {
		return nil // some err?
	}

	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: port})
	if err != nil {
		return err
	}
	g.listener = listener
	g.running = true

	go g.callback(listener)
	log.Println("listen udp port ", port, ", server start success")
	return nil
}

func (g *GeekerNetUDPServer) parseCmd(addr *net.UDPAddr, cmd string) {
	fmt.Println("get command ", cmd)
	list := strings.SplitN(cmd, "#", 2)
	g.cmdMap[list[0]](addr, list[1])
}

func (g *GeekerNetUDPServer) closeAll() {
	g.listener.Close()
	g.listener = nil
}

// gnetUDPServerCallback server callback.
func (g *GeekerNetUDPServer) callback(listener *net.UDPConn) {
	defer g.closeAll()
	data := make([]byte, 1024)
	for g.running {
		n, udpAddr, err := listener.ReadFromUDP(data)
		if err != nil {
			log.Println("error when read : ", err)
			continue
		}

		g.parseCmd(udpAddr, string(data[:n]))
	}

}

// Search search a seesionId's info.
func (g *GeekerNetUDPServer) Search(sessionID string) GeekerNodeInfo {
	return g.sessionMgr.Search(sessionID)
}

// SetSessionManager set session manager(cause not design very well)
func (g *GeekerNetUDPServer) SetSessionManager(mgr ISessionManager) {
	g.sessionMgr = mgr
}

// sessionManager session manager(not consiter session timeout at all).
type sessionManager struct {
	sessionMap map[string]GeekerNodeInfo
}

// Register register a session.
func (s *sessionManager) Register(sessionID string, info GeekerNodeInfo) error {
	s.sessionMap[sessionID] = info
	return nil
}

// Search search a sessionId.
func (s *sessionManager) Search(sessionID string) GeekerNodeInfo {
	return s.sessionMap[sessionID]
}
