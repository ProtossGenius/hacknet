package hacknetitf

import (
	"fmt"
	"net"

	"github.com/ProtossGenius/hacknet/hnlog"
	"github.com/ProtossGenius/hacknet/pb/cs"
	"github.com/ProtossGenius/hacknet/pb/hmsg"
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
	p2pHelper    pinfo.P2PHelperItf
	binder       *net.UDPConn
}

// AcceptHacker Hacker ask for login. extraData is hacker's detail info.
func (s *s4cImpl) AcceptHacker(udpAddr *net.UDPAddr, email string, extraData string) (result string) {
	point := s.pointInfoMgr.HackerJoin(udpAddr, email, extraData)

	return pinfo.GetHackerStatusName(point.Status)
}

// Hack connect to another Hacker's computer.
func (s *s4cImpl) Hack(hacker *pinfo.PointInfo, targetEmail string, extraData string) (result string) {
	targetHacker := s.pointInfoMgr.FindHacker(targetEmail)
	if targetHacker == nil || targetHacker.Status != pinfo.HackerStatusLive {
		return ErrNoSuchHacker
	}

	room, err := s.p2pHelper.CreateRoom(hacker, targetHacker)
	if err != nil {
		return err.Error()
	}

	go s.help2p(room)

	return ""
}

func (s *s4cImpl) help2p(room *pinfo.P2PRoom) {
}

// details log info's details.
type details map[string]interface{}

// MaxPackageSize udp package's max size.
const MaxPackageSize = 1024

func (s *s4cImpl) dealPackage(msg *hmsg.Message, hackerAddr *net.UDPAddr,
	hackerInfo *pinfo.PointInfo) (string, details, error) {
	wrapError := func(err error) error {
		return fmt.Errorf("s4cImpl.dealPackage, Error : %w", err)
	}

	hnlog.Info("accept data", logrus.Fields{"remoteAddr": hackerAddr, "message": msg, "hackerInfo": hackerInfo})

	if msg.Enum != int32(smn_dict.EDict_cs_Register) && hackerInfo == nil {
		return "check hackerInfo, not exist", details{"email": msg.Email}, errors.New(ErrNoSuchHacker)
	}

	switch msg.Enum {
	case int32(smn_dict.EDict_cs_Register):
		registerMsg := new(cs.Register)
		if err := proto.Unmarshal([]byte(msg.Msg), registerMsg); err != nil {
			return "unmarshal msg.Msg", details{"msg.Enum": msg.Enum, "msg.Msg": msg.Msg}, wrapError(err)
		}

		s.AcceptHacker(hackerAddr, registerMsg.Email, registerMsg.PubKey)
	case int32(smn_dict.EDict_cs_AskHack):
		askHackMsg := new(cs.AskHack)
		if err := proto.Unmarshal([]byte(msg.Msg), askHackMsg); err != nil {
			return "unmarshal msg.Msg", details{"msg.Enum": msg.Enum, "msg.Msg": msg.Msg}, wrapError(err)
		}

	default:
		return "unknow Enum", details{"msg.Enum": msg.Enum, "msg.Msg": msg.Msg}, wrapError(errors.New(ErrUnexceptEnum))
	}

	return "", nil, nil
}

func (s *s4cImpl) readMsg(data []byte) (*hmsg.Message, *pinfo.PointInfo, error) {
	msg := new(hmsg.Message)
	if err := proto.Unmarshal(data, msg); err != nil {
		return nil, nil, fmt.Errorf("s4cImpl.readMsg error : %w", err)
	}

	if hacker := s.pointInfoMgr.FindHacker(msg.Email); hacker != nil {
		return msg, hacker, nil
	}

	return msg, nil, nil
}

func (s *s4cImpl) startUp() {
	bytes := make([]byte, MaxPackageSize)

	for {
		readSize, hackerAddr, err := s.binder.ReadFromUDP(bytes)
		if err != nil {
			hnlog.Error("s4cImpl.startUp#readMsg", details{"error": err})

			continue
		}

		go func() {
			data := make([]byte, readSize)
			copy(data, bytes[0:readSize])

			if msg, hacker, err := s.readMsg(data); err != nil {
				hnlog.Error("s4cImpl.startUp#readMsg", details{"error": err})
			} else if info, fields, err := s.dealPackage(msg, hackerAddr, hacker); err != nil {
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
	res := &s4cImpl{pointInfoMgr: pinfo.NewPointInfoMgr(), binder: nil, p2pHelper: pinfo.NewP2PHelper()}
	// udp bind port
	var err error
	res.binder, err = net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: port, Zone: ""})

	check(err)

	go res.startUp()

	return res
}
