package hackeritf

import (
	"net"

	"github.com/ProtossGenius/hacknet/hacknetitf"
	"github.com/ProtossGenius/hacknet/hnlog"
	"github.com/ProtossGenius/hacknet/pb/hnp"
	"github.com/ProtossGenius/hacknet/pb/smn_dict"
)

func onForwardMsg(email string, serverAddr *net.UDPAddr, msg *hnp.ForwardMsg) (
	string, map[string]interface{}, error) {
	hnlog.Info("get forwardMsg", hnlog.Fields{"email": email, "serverAddr": serverAddr, "msg": msg})

	switch msg.Enums {
	case int32(smn_dict.EDict_hnep_HackAsk):
		{
		}

	case int32(smn_dict.EDict_hnep_BeHackAns):
	default:
		hnlog.Info("undeal msg", hnlog.Fields{"msg": msg})
	}

	return "", nil, nil
}

func onResult(email string, serverAddr *net.UDPAddr, msg *hnp.Result) (
	string, map[string]interface{}, error) {
	hnlog.Info("get resultMsg", hnlog.Fields{"email": email, "serverAddr": serverAddr, "msg": msg})

	if msg.Enums == int32(smn_dict.EDict_hnep_StrMsg) {
		// anything.
	}

	return "", nil, nil
}

// memHacker impl HackerItf.
type memHacker struct {
	server hacknetitf.ServerItf
}

// DoHack connect two port.
func (m *memHacker) DoHack(localPort int, targetEmail string, targetPort int) {
}

// GetOnForwardMsg get on forward msg for seerverItf use.
func (m *memHacker) GetOnForwardMsg() hacknetitf.OnForwardMsg {
	return onForwardMsg
}

// GetOnResultMsg get on result msg for serverItf use.
func (m *memHacker) GetOnResultMsg() hacknetitf.OnResultMsg {
	return onResult
}

// SetServer set local server Itf.
func (m *memHacker) SetServer(server hacknetitf.ServerItf) {
	m.server = server
}

func newMemHacker() HackerItf {
	res := &memHacker{server: nil}

	return res
}
