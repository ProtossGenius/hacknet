package hacknetitf

import (
	"net"

	"github.com/ProtossGenius/hacknet/pinfo"
)

// ServerForClientItf server for client.
type ServerForClientItf interface {
	// AcceptHacker Hacker ask for login. extraData is hacker's detail info.
	AcceptHacker(hackerAddr *net.UDPAddr, email, extraData string) (result string)
	// Hack connect to another Hacker's computer.
	Hack(hacker *pinfo.PointInfo, targetEmail, extraData string) (result string)
}

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
	// MoreConnect create more connect.
	MoreConnect() (ip string, port int, err string)
}
