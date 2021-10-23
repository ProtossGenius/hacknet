package udpfwd

import (
	"net"

	"github.com/ProtossGenius/hacknet/hnlog"
	"github.com/ProtossGenius/hacknet/pb/udpfwd"
	"google.golang.org/protobuf/proto"
)

// UDPForward forward udp msg.
type UDPForward struct {
	conn *net.UDPConn
	addr *net.UDPAddr
}

// NewUDPForward new .
func NewUDPForward(conn *net.UDPConn, addr *net.UDPAddr) *UDPForward {
	return &UDPForward{conn: conn, addr: addr}
}

func (u *UDPForward) sendMsg(addr *net.UDPAddr, data []byte) {
	if addr.IP.Equal(u.addr.IP) && addr.Port == u.addr.Port { // from another point.
		msg := new(udpfwd.UDPFwdMsg)
		if err := proto.Unmarshal(data, msg); err != nil {
			hnlog.Error("unmarshal error", hnlog.Fields{"err": err})

			return
		}

		if size, err := u.conn.WriteToUDP(msg.Msg,
			&net.UDPAddr{
				IP:   net.ParseIP(msg.IP),
				Port: int(msg.Port),
				Zone: "",
			}); err != nil || size != len(msg.Msg) {
			hnlog.Error("send Msg error", hnlog.Fields{"err": err})

			return
		}
	} else {
		msg := &udpfwd.UDPFwdMsg{Msg: data, IP: addr.IP.String(), Port: int32(addr.Port)}
		if size, err := u.conn.WriteToUDP(msg.Msg, u.addr); err != nil || size != len(msg.Msg) {
			hnlog.Error("send Msg error", hnlog.Fields{"err": err})

			return
		}
	}
}

// Work listen and forward.
func (u *UDPForward) Work() {
	const MaxSize = 2000
	bytes := make([]byte, MaxSize)

	for {
		size, addr, err := u.conn.ReadFromUDP(bytes)
		if err != nil {
			hnlog.Error("read UDP error", hnlog.Fields{"err": err})

			return
		}

		if size <= 0 {
			continue
		}

		data := make([]byte, size)
		copy(data, bytes[0:size])

		go u.sendMsg(addr, data)
	}
}
