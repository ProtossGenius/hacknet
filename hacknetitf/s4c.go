package hacknetitf

import (
	"net"

	"github.com/ProtossGenius/hacknet/pinfo"
)

// s4cImpl ServerForClientItf's impl.
type s4cImpl struct {
	pointInfoMgr pinfo.PointInfoMgrItf
	binder       *net.UDPConn
}

// AcceptHacker Hacker ask for login. extraData is hacker's detail info.
func (s *s4cImpl) AcceptHacker(udpAddr *net.UDPAddr, email string, extraData string) (result string) {
	point := s.pointInfoMgr.HackerJoin(udpAddr, email, extraData)

	return point.StatusJson()
}

// Hack connect to another Hacker's computer.
func (s *s4cImpl) Hack(hackerEmail string, hackerHost string, targetEmail string, extraData string) (result string) {
	panic("not implemented") // TODO: Implement
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// news4c new s4cImpl.
func news4c(port int) ServerForClientItf {
	res := &s4cImpl{pointInfoMgr: pinfo.NewPointInfoMgr(), binder: nil}
	// udp bind port
	var err error
	res.binder, err = net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: port, Zone: ""})

	check(err)

	return res
}
