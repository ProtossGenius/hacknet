package hacknetitf

import (
	"fmt"
	"net"

	"github.com/ProtossGenius/hacknet/hnlog"
	"github.com/ProtossGenius/hacknet/pb/hmsg"
	"github.com/ProtossGenius/hacknet/pb/sc"
	"github.com/ProtossGenius/hacknet/pb/smn_dict"
	"github.com/ProtossGenius/hacknet/pinfo"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/internal/errors"
	"google.golang.org/protobuf/proto"
)

const (
	// ErrNoSuchHacker can't found this hacker by email.
	ErrNoSuchHacker = "ErrNoSuchHacker"
	// ErrUnexceptEnum enum not in except in this deal function.
	ErrUnexceptEnum = "ErrUnexceptEnum"
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

// details log info's details.
type details map[string]interface{}

// MaxPackageSize udp package's max size.
const MaxPackageSize = 1024

func (s *s4cImpl) dealPackage(data []byte, remoteAddr *net.UDPAddr, err error) (string, details, error) {
	if err != nil {
		return "accept package.", details{}, err
	}

	wrapError := func(err error) error {
		return fmt.Errorf("s4cImpl.dealPackage, Error : %w", err)
	}

	hnlog.Info("accept data", logrus.Fields{"remoteAddr": remoteAddr, "readSize": len(data), "message": data})

	msg := new(hmsg.Message)
	// TODO: decode to protoc and do something.
	err = proto.Unmarshal(data, msg)
	if err != nil {
		return "unmarshal msg", details{}, wrapError(err)
	}

	switch msg.Enum {
	case int32(smn_dict.EDict_cs_Register):
		registerMsg := new(sc.AnsHack)
		err = proto.Unmarshal([]byte(msg.Msg), registerMsg)

		if err != nil {
			return "unmarshal msg.Msg", details{"msg.Enum": smn_dict.EDict_sc_AnsHack, "msg.Msg": msg.Msg}, wrapError(err)
		}

	default:
		return "unknow Enum", details{"msg.Enum": msg.Enum, "msg.Msg": msg.Msg}, wrapError(errors.New(ErrUnexceptEnum))
	}

	return "", nil, nil
}

func (s *s4cImpl) startUp() {
	bytes := make([]byte, MaxPackageSize)

	for {
		n, remoteAddr, err := s.binder.ReadFromUDP(bytes)

		go func() {
			data := make([]byte, n)
			copy(data, bytes[0:n])

			if info, fields, err := s.dealPackage(data, remoteAddr, err); err != nil {
				if fields == nil {
					fields = details{}
					fields["fields_not_init"] = true
				}

				fields["error"] = err
				hnlog.Error(info, fields)
			}
		}()
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
