package hacknetitf

import (
	"net"

	"github.com/ProtossGenius/hacknet/hnlog"
	"github.com/ProtossGenius/hacknet/pb/hmsg"
	"github.com/ProtossGenius/hacknet/pb/hnp"
	"github.com/ProtossGenius/hacknet/pinfo"
	"google.golang.org/protobuf/proto"
)

/*@SMIST
include("parseProtos.js");
setIgnoreInput(true);
proto2GoItf("./datas/protos/hnp.proto", "ServerItf", "server for client")
*/
// ServerItf server for client.
type ServerItf interface {
	// GetUDPConn get server's udp conn.
	GetUDPConn() *net.UDPConn
	// GetPointMgr get server's poingMgr.
	GetPointMgr() pinfo.PointInfoMgrItf
	// Register register this client to server.
	Register(email string, hackerAddr *net.UDPAddr, msg *hnp.Register) (string, map[string]interface{}, error)
	// Result register result
	Result(email string, hackerAddr *net.UDPAddr, msg *hnp.Result) (string, map[string]interface{}, error)
	// CheckEmail check if email belong to register.
	CheckEmail(email string, hackerAddr *net.UDPAddr, msg *hnp.CheckEmail) (string, map[string]interface{}, error)
	// Forward send email to another point.
	Forward(email string, hackerAddr *net.UDPAddr, msg *hnp.Forward) (string, map[string]interface{}, error)
	// ForwardMsg send message to the point.
	ForwardMsg(email string, hackerAddr *net.UDPAddr, msg *hnp.ForwardMsg) (string, map[string]interface{}, error)
	// HeartJump heart jump just for keep alive.
	HeartJump(email string, hackerAddr *net.UDPAddr, msg *hnp.HeartJump) (string, map[string]interface{}, error)
}

// @SMIST setIgnoreInput(false)

// ServerFactory product server for client.
type ServerFactory func(port int, email, pubKey string, onForwardMsg OnForwardMsg,
	onResultMsg OnResultMsg) ServerItf

// SetServerForClientFactory set factory.
func SetServerForClientFactory(factory ServerFactory) {
	if factory != nil {
		serverForClientFactory = factory
	}
}

// NewServer create.
func NewServer(port int, email, pubKey string, onForwardMsg OnForwardMsg, onResultMsg OnResultMsg) ServerItf {
	return serverForClientFactory(port, email, pubKey, onForwardMsg, onResultMsg)
}

// writeMsg write message to binder.
func writeMsg(binder *net.UDPConn, hackerAddr *net.UDPAddr, msg *hmsg.Message) {
	var err error

	var data []byte

	if data, err = proto.Marshal(msg); err != nil {
		hnlog.Error("s4cImpl.write Marshal", details{"msg": msg, "err": err})

		return
	}

	sendSize, err := binder.WriteToUDP(data, hackerAddr)
	if err != nil || sendSize != len(data) {
		hnlog.Error("s4cImpl.write send", details{"sendSize": sendSize, "err": err, "data": data, "msg": msg})
	}
}
