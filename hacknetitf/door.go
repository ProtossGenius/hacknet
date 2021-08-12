package hacknetitf

import (
	"net"

	"github.com/ProtossGenius/hacknet/hnlog"
	"github.com/ProtossGenius/hacknet/pb/cs"
	"github.com/ProtossGenius/hacknet/pb/hmsg"
	"google.golang.org/protobuf/proto"
)

/*@SMIST
include("parseProtos.js");
setIgnoreInput(true);
proto2GoItf("./protos/cs.proto", "ServerItf", "server for client")
*/
// ServerItf server for client.
type ServerItf interface {
	// Register register this client to server.
	Register(email string, hackerAddr *net.UDPAddr, msg *cs.Register) (*hmsg.Message, map[string]interface{}, error)
	// CheckEmail check if email belong to register.
	CheckEmail(email string, hackerAddr *net.UDPAddr, msg *cs.CheckEmail) (*hmsg.Message, map[string]interface{}, error)
	// Forward send email to another point.
	Forward(email string, hackerAddr *net.UDPAddr, msg *cs.Forward) (*hmsg.Message, map[string]interface{}, error)
	// SendMsg send message to the point.
	SendMsg(email string, hackerAddr *net.UDPAddr, msg *cs.SendMsg) (*hmsg.Message, map[string]interface{}, error)
	// HeartJump heart jump just for keep alive.
	HeartJump(email string, hackerAddr *net.UDPAddr, msg *cs.HeartJump) (*hmsg.Message, map[string]interface{}, error)
}
// @SMIST setIgnoreInput(false)

// ServerForClientFactory product server for client.
type ServerForClientFactory func(port int) ServerItf

// SetServerForClientFactory set factory.
func SetServerForClientFactory(factory ServerForClientFactory) {
	if factory != nil {
		serverForClientFactory = factory
	}
}

// NewServerForClient create.
func NewServerForClient(port int) ServerItf {
	return serverForClientFactory(port)
}

// ClientForClientItf client giver service to another client(client are connected).
type ClientForClientItf interface {
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

// addrEquals is addr equals.
func addrEquals(lhs, rhs *net.UDPAddr) bool {
	return lhs.IP.Equal(rhs.IP) && lhs.Port == rhs.Port
}
