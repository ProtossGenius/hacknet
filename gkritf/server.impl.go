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
			infoMap:    make(map[GeekerNodeInfo]string),
		},
		running: false,
	}).init()
}

// FCommand what will command do.
type FCommand func(addr *net.UDPAddr, msg GeekerMsg) error

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

	sendmsg := func(addr *net.UDPAddr, msg GeekerMsg) (err error) {
		if _, err = g.listener.WriteToUDP([]byte(msg), addr); err != nil {
			log.Println("when send response msg : ", msg, ", error is ", err)
		}

		log.Println("send msg to ", addr, ":", msg)

		return err
	}

	g.cmdMap = map[string]FCommand{
		"SearchC": func(addr *net.UDPAddr, info GeekerMsg) (err error) {
			if searchC, err := NewSearchC(info); err == nil {
				return sendmsg(addr, SearchS{g.Search(searchC.SessionID)}.Message())
			}

			return err

		},
		"RegisterC": func(addr *net.UDPAddr, msg GeekerMsg) (err error) {
			var registerC RegisterC
			if registerC, err = NewRegisterC(msg); err != nil {
				return err
			}
			g.sessionMgr.Register(registerC.SessionID, GeekerNodeInfo{IP: addr.IP.String(), Port: addr.Port})

			return sendmsg(addr, "Register:done")
		},
		"NoticeC": func(addr *net.UDPAddr, noticeInfo GeekerMsg) (err error) {
			var noticeC NoticeC
			if noticeC, err = NewNoticeC(noticeInfo); err != nil {
				return err
			}

			targetInfo := g.Search(noticeC.TargetSession)
			nodeInfo := GeekerNodeInfo{IP: addr.IP.String(), Port: addr.Port}
			if targetInfo.Inited() {
				sendmsg(targetInfo.BuildUDPAddr(), NoticeS{
					NodeSessionID: g.sessionMgr.SearchUUID(nodeInfo),
					NodeInfo:      nodeInfo,
					ExtraData:     noticeC.ExtraData,
				}.Message())
			}
			return sendmsg(addr, "Notice#"+noticeInfo+":done")
		},
		"NetInfoC": func(addr *net.UDPAddr, msg GeekerMsg) (err error) {
			return sendmsg(addr, NetInfoS{IP: addr.IP.String(), Port: addr.Port}.Message())
		},
	}

	return g
}

// Listen listion a local udp port.
func (g *GeekerNetUDPServer) Listen(port int) error {
	if g.running {
		return nil // some err?
	}

	g.sessionMgr.Startup()

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
	if len(list) != 2 {
		return
	}
	g.cmdMap[list[0]](addr, GeekerMsg(list[1]))
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

// sessionPair session info pair.
type sessionPair struct {
	sessionID string
	info      GeekerNodeInfo
}

// sessionManager session manager(not consiter session timeout at all).
type sessionManager struct {
	sessionMap map[string]GeekerNodeInfo
	infoMap    map[GeekerNodeInfo]string
	put        chan sessionPair
	del        chan string
	getInfo    chan string
	infoResult chan GeekerNodeInfo
	getSSID    chan GeekerNodeInfo
	ssidResult chan string
	running    bool
}

func (s *sessionManager) SearchUUID(info GeekerNodeInfo) string {
	s.getSSID <- info
	return <-s.ssidResult
}

func (s *sessionManager) Close() {
	s.running = false
	close(s.put)
	close(s.del)
	close(s.getInfo)
	close(s.infoResult)
	close(s.getSSID)
	close(s.ssidResult)
}

// Startup start up session manager.
func (s *sessionManager) Startup() error {
	s.put = make(chan sessionPair, 1000)
	s.del = make(chan string, 1000)
	s.getInfo = make(chan string, 1000)
	s.infoResult = make(chan GeekerNodeInfo, 1000)
	s.getSSID = make(chan GeekerNodeInfo, 1000)
	s.ssidResult = make(chan string, 1000)
	s.running = true
	go func() {
		for s.running {
			select {
			case pair := <-s.put:
				s.sessionMap[pair.sessionID] = pair.info
				s.infoMap[pair.info] = pair.sessionID
			case todel := <-s.del:
				tdInfo := s.sessionMap[todel]
				delete(s.sessionMap, todel)
				delete(s.infoMap, tdInfo)
			case ssid := <-s.getInfo:
				s.infoResult <- s.sessionMap[ssid]
			case info := <-s.getSSID:
				s.ssidResult <- s.infoMap[info]
			}
		}
	}()

	return nil
}

// Register register a session.
func (s *sessionManager) Register(sessionID string, info GeekerNodeInfo) error {
	s.put <- sessionPair{sessionID: sessionID, info: info}
	if len(s.sessionMap) > 20000 {
		times := 30
		for k := range s.sessionMap {
			s.del <- k
			times--
			if times <= 0 {
				return nil
			}
		}
	}
	return nil
}

// Search search a sessionId.
func (s *sessionManager) Search(sessionID string) GeekerNodeInfo {
	s.getInfo <- sessionID
	return <-s.infoResult
}
