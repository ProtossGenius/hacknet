package hacknetitf

import (
	"net"

	"github.com/ProtossGenius/hacknet/pb/cs"
	"google.golang.org/protobuf/proto"
)

/*@SMIST
include("parseProtos.js");
setIgnoreInput(true);
proto2GoItf("./protos/cs.proto", "ServerForClientItf", "server for client")
*/
// ServerForClientItf server for client.
type ServerForClientItf interface {
	// Register register this client to server.
	Register(email string, hackerAddr *net.UDPAddr, msg *cs.Register) (*proto.Message, map[string]interface{}, error)
	// CheckEmail check if email belong to register.
	CheckEmail(email string, hackerAddr *net.UDPAddr, msg *cs.CheckEmail) (*proto.Message, map[string]interface{}, error)
	// AskHack ask connect another client.
	AskHack(email string, hackerAddr *net.UDPAddr, msg *cs.AskHack) (*proto.Message, map[string]interface{}, error)
	// HeartJump heart jump just for keep alive.
	HeartJump(email string, hackerAddr *net.UDPAddr, msg *cs.HeartJump) (*proto.Message, map[string]interface{}, error)
}
// @SMIST setIgnoreInput(false)

// ServerForClientFactory product server for client.
type ServerForClientFactory func(port int) ServerForClientItf

// SetServerForClientFactory set factory.
func SetServerForClientFactory(factory ServerForClientFactory) {
	if factory != nil {
		serverForClientFactory = factory
	}
}

// NewServerForClient create.
func NewServerForClient(port int) ServerForClientItf {
	return serverForClientFactory(port)
}

// ClientForClientItf client giver service to another client(client are connected).
type ClientForClientItf interface {
}
