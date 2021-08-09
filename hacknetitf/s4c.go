package hacknetitf

import (
	"net"

	"github.com/ProtossGenius/hacknet/hnlog"
	"github.com/ProtossGenius/hacknet/pinfo"
	"github.com/sirupsen/logrus"
)

const (
	// ErrNoSuchHacker can't found this hacker by email.
	ErrNoSuchHacker = "ErrNoSuchHacker"
)

// s4cImpl ServerForClientItf's impl.
type s4cImpl struct {
	pointInfoMgr pinfo.PointInfoMgrItf
	binder       *net.UDPConn
}

// AcceptHacker Hacker ask for login. extraData is hacker's detail info.
func (s *s4cImpl) AcceptHacker(udpAddr *net.UDPAddr, email string, extraData string) (result string) {
	point := s.pointInfoMgr.HackerJoin(udpAddr, email, extraData)

	return pinfo.GetHackerStatusName(point.Status)
}

// Hack connect to another Hacker's computer.
func (s *s4cImpl) Hack(hackerEmail string, hackerAddr *net.UDPAddr, targetEmail string,
	extraData string) (result string) {
	hacker := s.pointInfoMgr.FindHacker(hackerEmail)
	if hacker == nil {
		return ErrNoSuchHacker
	}

	return ""
}

func (s *s4cImpl) startUp() {
	bytes := make([]byte, 1024)
	for {
		n, remoteAddr, err := s.binder.ReadFromUDP(bytes)
		if err != nil {
			hnlog.Error(err.Error(), logrus.Fields{})
			continue
		}
		hnlog.Info("accept data", logrus.Fields{"remoteAddr": remoteAddr, "readSize": n, "message": bytes[0:n]})

		// TODO: decode to protoc and do something.
	}
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

	go res.startUp()

	return res
}
