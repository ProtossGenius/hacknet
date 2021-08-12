package hacknetitf

import (
	"errors"
	"fmt"
	"net"

	"github.com/ProtossGenius/hacknet/hnlog"
	"github.com/ProtossGenius/hacknet/pb/cs"
	"github.com/ProtossGenius/hacknet/pb/hmsg"
	"github.com/ProtossGenius/hacknet/pb/smn_dict"
	"github.com/ProtossGenius/hacknet/pinfo"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

const (
	// ErrNoSuchHacker can't found this hacker by email.
	ErrNoSuchHacker = "ErrNoSuchHacker"
	// ErrUnexceptEnum enum not in except in this deal function.
	ErrUnexceptEnum = "ErrUnexceptEnum"
)

// ErrHackerExist hacker is exist.
var ErrHackerExist = errors.New("ErrHackerExist")

// s4cImpl ServerItf's impl.
type s4cImpl struct {
	pointInfoMgr pinfo.PointInfoMgrItf
	p2pHelper    pinfo.P2PHelperItf
	binder       *net.UDPConn
}

// Register register this client to server.
func (s *s4cImpl) Register(email string, hackerAddr *net.UDPAddr, msg *cs.Register) (
	*hmsg.Message, map[string]interface{}, error) {
	findH := s.pointInfoMgr.FindHacker(email)
	if findH != nil && !addrEquals(hackerAddr, findH.HackerAddr) {
		return nil, details{}, ErrHackerExist
	}

	_ = s.pointInfoMgr.HackerJoin(hackerAddr, email, msg.PubKey)

	return &hmsg.Message{Enum: int32(smn_dict.EDict_cs_SendMsg), Email: email, Msg: ""}, details{}, nil
}

// CheckEmail check if email belong to register.
func (s *s4cImpl) CheckEmail(email string, hackerAddr *net.UDPAddr, msg *cs.CheckEmail) (
	*hmsg.Message, map[string]interface{}, error) {
	panic("not implemented") // TODO: Implement
}

// Forward send email to another point.
func (s *s4cImpl) Forward(email string, hackerAddr *net.UDPAddr, msg *cs.Forward) (
	*hmsg.Message, map[string]interface{}, error) {
	panic("not implemented") // TODO: Implement
}

// SendMsg send message to the point.
func (s *s4cImpl) SendMsg(email string, hackerAddr *net.UDPAddr, msg *cs.SendMsg) (
	*hmsg.Message, map[string]interface{}, error) {
	panic("not implemented") // TODO: Implement
}

// HeartJump heart jump just for keep alive.
func (s *s4cImpl) HeartJump(email string, hackerAddr *net.UDPAddr, msg *cs.HeartJump) (
	*hmsg.Message, map[string]interface{}, error) {
	panic("not implemented") // TODO: Implement
}

// details log info's details.
type details map[string]interface{}

// MaxPackageSize udp package's max size.
const MaxPackageSize = 1024

// dealPackage deal package.
func (s *s4cImpl) dealPackage(msg *hmsg.Message, hackerAddr *net.UDPAddr,
	hackerInfo *pinfo.PointInfo) (string, details, error) {
	binder := s.binder
	wrapError := func(err error) error {
		return fmt.Errorf("s4cImpl.dealPackage, Error : %w", err)
	}

	var err error

	hnlog.Info("accept data", logrus.Fields{"remoteAddr": hackerAddr, "message": msg, "hackerInfo": hackerInfo})

	if msg.Enum != int32(smn_dict.EDict_cs_Register) && hackerInfo == nil {
		return "check hackerInfo, not exist", details{"email": msg.Email}, wrapError(errors.New(ErrNoSuchHacker))
	}

	// @SMIST include("parseProtos.js"); proto2GoSwitch("./protos/cs.proto", 1)
	const UnmarshalMsgMsg = "Unmarshal msg.Msg"

	var _resp *hmsg.Message

	var detail details

	switch msg.Enum {
	case int32(smn_dict.EDict_cs_Register):
		_subMsg := new(cs.Register)
		if err = proto.Unmarshal([]byte(msg.Msg), _subMsg); err != nil {
			return UnmarshalMsgMsg, details{"msg.Enum": msg.Enum, "msg.Msg": msg.Msg}, wrapError(err)
		}

		if _resp, detail, err = s.Register(msg.Email, hackerAddr, _subMsg); err != nil {
			return "s.Register", detail, wrapError(err)
		}

		writeMsg(binder, hackerAddr, _resp)
	case int32(smn_dict.EDict_cs_CheckEmail):
		_subMsg := new(cs.CheckEmail)
		if err = proto.Unmarshal([]byte(msg.Msg), _subMsg); err != nil {
			return UnmarshalMsgMsg, details{"msg.Enum": msg.Enum, "msg.Msg": msg.Msg}, wrapError(err)
		}

		if _resp, detail, err = s.CheckEmail(msg.Email, hackerAddr, _subMsg); err != nil {
			return "s.CheckEmail", detail, wrapError(err)
		}

		writeMsg(binder, hackerAddr, _resp)
	case int32(smn_dict.EDict_cs_Forward):
		_subMsg := new(cs.Forward)
		if err = proto.Unmarshal([]byte(msg.Msg), _subMsg); err != nil {
			return UnmarshalMsgMsg, details{"msg.Enum": msg.Enum, "msg.Msg": msg.Msg}, wrapError(err)
		}

		if _resp, detail, err = s.Forward(msg.Email, hackerAddr, _subMsg); err != nil {
			return "s.Forward", detail, wrapError(err)
		}

		writeMsg(binder, hackerAddr, _resp)
	case int32(smn_dict.EDict_cs_SendMsg):
		_subMsg := new(cs.SendMsg)
		if err = proto.Unmarshal([]byte(msg.Msg), _subMsg); err != nil {
			return UnmarshalMsgMsg, details{"msg.Enum": msg.Enum, "msg.Msg": msg.Msg}, wrapError(err)
		}

		if _resp, detail, err = s.SendMsg(msg.Email, hackerAddr, _subMsg); err != nil {
			return "s.SendMsg", detail, wrapError(err)
		}

		writeMsg(binder, hackerAddr, _resp)
	case int32(smn_dict.EDict_cs_HeartJump):
		_subMsg := new(cs.HeartJump)
		if err = proto.Unmarshal([]byte(msg.Msg), _subMsg); err != nil {
			return UnmarshalMsgMsg, details{"msg.Enum": msg.Enum, "msg.Msg": msg.Msg}, wrapError(err)
		}

		if _resp, detail, err = s.HeartJump(msg.Email, hackerAddr, _subMsg); err != nil {
			return "s.HeartJump", detail, wrapError(err)
		}

		writeMsg(binder, hackerAddr, _resp)
	default:
		return "unknow Enum", details{"msg.Enum": msg.Enum, "msg.Msg": msg.Msg}, wrapError(errors.New(ErrUnexceptEnum))
	}
	/* @SMIST setIgnoreInput(false);*/
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
func news4c(port int) ServerItf {
	res := &s4cImpl{pointInfoMgr: pinfo.NewPointInfoMgr(), binder: nil, p2pHelper: pinfo.NewP2PHelper()}
	// udp bind port
	var err error
	res.binder, err = net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: port, Zone: ""})

	check(err)

	go res.startUp()

	return res
}
