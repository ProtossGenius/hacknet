package hackeritf

import (
	"github.com/ProtossGenius/hacknet/hacknetitf"
)

// HackerItf client interface, connect to another point and connect or listen.
type HackerItf interface {
	// DoHack connect two port .
	DoHack(localPort int, targetEmail string, targetPort int)
	// GetOnForwardMsg get on forward msg for seerverItf use.
	GetOnForwardMsg() hacknetitf.OnForwardMsg
	// GetOnResultMsg get on result msg for serverItf use.
	GetOnResultMsg() hacknetitf.OnResultMsg
	// SetServer set local server Itf.
	SetServer(server hacknetitf.ServerItf)
}

// HackerFactory hackerItf's factory.
type HackerFactory func() HackerItf

var hackerFactory = newMemHacker

// SetHackerFactory set hacker factory.
func SetHackerFactory(factory HackerFactory) {
	if factory != nil {
		hackerFactory = factory
	}
}

// NewHacker create new HackerItf.
func NewHacker() HackerItf {
	return hackerFactory()
}

// NewServerWithHacker .
func NewServerWithHacker(port int, email, pubKey string) (hacknetitf.ServerItf, HackerItf) {
	hacker := hackerFactory()

	return hacknetitf.NewServer(port, email, pubKey, hacker.GetOnForwardMsg(), hacker.GetOnResultMsg()),
		hacker
}
