package hacknetitf

import (
	"errors"
	"fmt"
	"net"

	"github.com/ProtossGenius/hacknet/hnlog"
	"github.com/ProtossGenius/hacknet/pb/hmsg"
	"github.com/ProtossGenius/hacknet/pb/hnp"
	"github.com/ProtossGenius/hacknet/pb/smn_dict"
	"github.com/ProtossGenius/hacknet/pinfo"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

var (
	// ErrNoSuchHacker can't found this hacker by email.
	ErrNoSuchHacker = errors.New("ErrNoSuchHacker")
	// ErrUnexceptEnum enum not in except in this deal function.
	ErrUnexceptEnum = errors.New("ErrUnexceptEnum")
	// ErrHackerExist hacker is exist.
	ErrHackerExist = errors.New("ErrHackerExist")
)

const (
	// FailAddrNotEqual ip not equal.
	FailAddrNotEqual = "FailIpNotEqual"
)

// s4cImpl ServerItf's impl.
type s4cImpl struct {
	pointInfoMgr pinfo.PointInfoMgrItf
	binder       *net.UDPConn
	email        string
	// callbacks
	onForwardMsg OnForwardMsg
	onResultMsg  OnResultMsg
	pubKey       string
}

// GetPointMgr get server's poingMgr.
func (s *s4cImpl) GetPointMgr() pinfo.PointInfoMgrItf {
	return s.pointInfoMgr
}

// GetUDPConn get server's udp conn.
func (s *s4cImpl) GetUDPConn() *net.UDPConn {
	return s.binder
}

// Register register this client to server.
func (s *s4cImpl) Register(email string, hackerAddr *net.UDPAddr, msg *hnp.Register) (
	string, map[string]interface{}, error) {
	var hackerInfo *pinfo.PointInfo
	if hackerInfo = s.pointInfoMgr.FindHacker(email); hackerInfo != nil {
		if hackerInfo.Email != email {
			return FailAddrNotEqual, details{}, nil
		}

		return "", nil, nil
	}

	hackerInfo = s.pointInfoMgr.HackerJoin(hackerAddr, email, msg.PubKey)

	hnlog.Info("s.Register", details{"email": email, "hackerAddr": hackerAddr, "pubKey": msg.PubKey})

	if regMsg, err := Pack_hnp_Register(s.email, &hnp.Register{PubKey: s.pubKey}); err == nil {
		writeMsg(s.binder, hackerAddr, regMsg)
	} else {
		hnlog.Error("s.Register#Pack_hnp_Register", details{"err": err})
	}

	return pinfo.GetHackerStatusName(hackerInfo.Status), details{}, nil
}

// Result register result.
func (s *s4cImpl) Result(email string, hackerAddr *net.UDPAddr, msg *hnp.Result) (
	string, map[string]interface{}, error) {
	return s.onResultMsg(email, hackerAddr, msg)
}

// CheckEmail check if email belong to register.
func (s *s4cImpl) CheckEmail(email string, hackerAddr *net.UDPAddr, msg *hnp.CheckEmail) (
	string, map[string]interface{}, error) {
	return "", nil, nil
}

// Forward send email to another point.
func (s *s4cImpl) Forward(email string, hackerAddr *net.UDPAddr, msg *hnp.Forward) (
	string, map[string]interface{}, error) {
	target := s.findHacker(msg.Target)
	if target == nil {
		return ErrNoSuchHacker.Error(), nil, nil
	}

	resp, err := Pack_hnp_ForwardMsg(s.email, &hnp.ForwardMsg{
		FromEmail: email, FromIp: hackerAddr.IP.String(),
		Msg: msg.Msg, FromPort: int32(hackerAddr.Port),
		Enums: msg.Enums,
	})
	if err != nil {
		return err.Error(), details{"function": "Forward", "when": "pack_hnp_ForwardMsg", "err": err}, err
	}

	writeMsg(s.binder, target.HackerAddr, resp)

	return "", nil, nil
}

// ForwardMsg send message to the point.
func (s *s4cImpl) ForwardMsg(email string, hackerAddr *net.UDPAddr, msg *hnp.ForwardMsg) (
	string, map[string]interface{}, error) {
	return s.onForwardMsg(email, hackerAddr, msg)
}

// OnForwardMsg on forward msg.
type OnForwardMsg func(email string, hackerAddr *net.UDPAddr, msg *hnp.ForwardMsg) (
	string, map[string]interface{}, error)

// OnResultMsg on result msg.
type OnResultMsg func(email string, hackerAddr *net.UDPAddr, msg *hnp.Result) (
	string, map[string]interface{}, error)

// HeartJump heart jump just for keep alive.
func (s *s4cImpl) HeartJump(email string, hackerAddr *net.UDPAddr, msg *hnp.HeartJump) (
	string, map[string]interface{}, error) {
	return "", nil, nil
}

func (s *s4cImpl) findHacker(email string) *pinfo.PointInfo {
	hacker := s.pointInfoMgr.FindHacker(email)
	if hacker == nil || hacker.Status != pinfo.HackerStatusLive {
		return nil
	}

	return hacker
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

	if msg.Enum != int32(smn_dict.EDict_hnp_Register) && hackerInfo == nil {
		return "check hackerInfo, not exist", details{"email": msg.Email, "msg.Enum": msg.Enum}, wrapError(ErrNoSuchHacker)
	}

	// @SMIST include("parseProtos.js"); proto2GoSwitch("./protos/hnp.proto", 1)
	const (
		PackResult      = "packResult"
		UnmarshalMsgMsg = "Unmarshal msg.Msg"
	)

	var _resp string

	var detail details

	var _result *hmsg.Message

	switch msg.Enum {
	case int32(smn_dict.EDict_hnp_Register):
		_subMsg := new(hnp.Register)
		if err = msg.Msg.UnmarshalTo(_subMsg); err != nil {
			return UnmarshalMsgMsg, details{"msg.Enum": msg.Enum, "msg.Msg": msg.Msg}, wrapError(err)
		}

		if _resp, detail, err = s.Register(msg.Email, hackerAddr, _subMsg); err != nil {
			return "s.Register", detail, wrapError(err)
		}

		hnlog.Info("dealPackage", details{"method": "s.Register", "_resp": _resp, "details": detail, "err": err})

		if _result, err = Pack_hnp_Result(msg.Email, &hnp.Result{
			Enums: int32(smn_dict.EDict_hnp_Register),
			Info:  _resp,
		}); err != nil {
			return PackResult, details{"email": msg.Email, "_resp": _resp, "error": err}, wrapError(err)
		}

		writeMsg(binder, hackerAddr, _result)
	case int32(smn_dict.EDict_hnp_Result):
		_subMsg := new(hnp.Result)
		if err = msg.Msg.UnmarshalTo(_subMsg); err != nil {
			return UnmarshalMsgMsg, details{"msg.Enum": msg.Enum, "msg.Msg": msg.Msg}, wrapError(err)
		}

		if _resp, detail, err = s.Result(msg.Email, hackerAddr, _subMsg); err != nil {
			return "s.Result", detail, wrapError(err)
		}

		hnlog.Info("dealPackage", details{"method": "s.Result", "_resp": _resp, "details": detail, "err": err})

	case int32(smn_dict.EDict_hnp_CheckEmail):
		_subMsg := new(hnp.CheckEmail)
		if err = msg.Msg.UnmarshalTo(_subMsg); err != nil {
			return UnmarshalMsgMsg, details{"msg.Enum": msg.Enum, "msg.Msg": msg.Msg}, wrapError(err)
		}

		if _resp, detail, err = s.CheckEmail(msg.Email, hackerAddr, _subMsg); err != nil {
			return "s.CheckEmail", detail, wrapError(err)
		}

		hnlog.Info("dealPackage", details{"method": "s.CheckEmail", "_resp": _resp, "details": detail, "err": err})

		if _result, err = Pack_hnp_Result(msg.Email, &hnp.Result{
			Enums: int32(smn_dict.EDict_hnp_CheckEmail),
			Info:  _resp,
		}); err != nil {
			return PackResult, details{"email": msg.Email, "_resp": _resp, "error": err}, wrapError(err)
		}

		writeMsg(binder, hackerAddr, _result)
	case int32(smn_dict.EDict_hnp_Forward):
		_subMsg := new(hnp.Forward)
		if err = msg.Msg.UnmarshalTo(_subMsg); err != nil {
			return UnmarshalMsgMsg, details{"msg.Enum": msg.Enum, "msg.Msg": msg.Msg}, wrapError(err)
		}

		if _resp, detail, err = s.Forward(msg.Email, hackerAddr, _subMsg); err != nil {
			return "s.Forward", detail, wrapError(err)
		}

		hnlog.Info("dealPackage", details{"method": "s.Forward", "_resp": _resp, "details": detail, "err": err})

		if _result, err = Pack_hnp_Result(msg.Email, &hnp.Result{
			Enums: int32(smn_dict.EDict_hnp_Forward),
			Info:  _resp,
		}); err != nil {
			return PackResult, details{"email": msg.Email, "_resp": _resp, "error": err}, wrapError(err)
		}

		writeMsg(binder, hackerAddr, _result)
	case int32(smn_dict.EDict_hnp_ForwardMsg):
		_subMsg := new(hnp.ForwardMsg)
		if err = msg.Msg.UnmarshalTo(_subMsg); err != nil {
			return UnmarshalMsgMsg, details{"msg.Enum": msg.Enum, "msg.Msg": msg.Msg}, wrapError(err)
		}

		if _resp, detail, err = s.ForwardMsg(msg.Email, hackerAddr, _subMsg); err != nil {
			return "s.ForwardMsg", detail, wrapError(err)
		}

		hnlog.Info("dealPackage", details{"method": "s.ForwardMsg", "_resp": _resp, "details": detail, "err": err})

		if _result, err = Pack_hnp_Result(msg.Email, &hnp.Result{
			Enums: int32(smn_dict.EDict_hnp_ForwardMsg),
			Info:  _resp,
		}); err != nil {
			return PackResult, details{"email": msg.Email, "_resp": _resp, "error": err}, wrapError(err)
		}

		writeMsg(binder, hackerAddr, _result)
	case int32(smn_dict.EDict_hnp_HeartJump):
		_subMsg := new(hnp.HeartJump)
		if err = msg.Msg.UnmarshalTo(_subMsg); err != nil {
			return UnmarshalMsgMsg, details{"msg.Enum": msg.Enum, "msg.Msg": msg.Msg}, wrapError(err)
		}

		if _resp, detail, err = s.HeartJump(msg.Email, hackerAddr, _subMsg); err != nil {
			return "s.HeartJump", detail, wrapError(err)
		}

		hnlog.Info("dealPackage", details{"method": "s.HeartJump", "_resp": _resp, "details": detail, "err": err})

		if _result, err = Pack_hnp_Result(msg.Email, &hnp.Result{
			Enums: int32(smn_dict.EDict_hnp_HeartJump),
			Info:  _resp,
		}); err != nil {
			return PackResult, details{"email": msg.Email, "_resp": _resp, "error": err}, wrapError(err)
		}

		writeMsg(binder, hackerAddr, _result)
	default:
		return "unknow Enum", details{"msg.Enum": msg.Enum, "msg.Msg": msg.Msg}, wrapError(ErrUnexceptEnum)
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

		data := make([]byte, readSize)
		copy(data, bytes[0:readSize])

		go func() {
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
func news4c(port int, email, pubKey string, onForwardMsg OnForwardMsg, onResultMsg OnResultMsg) ServerItf {
	res := &s4cImpl{
		pointInfoMgr: pinfo.NewPointInfoMgr(), binder: nil,
		email: email, onForwardMsg: onForwardMsg, onResultMsg: onResultMsg, pubKey: pubKey,
	}

	// udp bind port
	var err error
	res.binder, err = net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: port, Zone: ""})

	check(err)

	go res.startUp()

	return res
}
