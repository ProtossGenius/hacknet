package main

import (
	"flag"
	"net"
	"time"

	"github.com/ProtossGenius/hacknet/hacknetitf"
	"github.com/ProtossGenius/hacknet/hnlog"
	"github.com/ProtossGenius/hacknet/pb/hnp"
	"github.com/ProtossGenius/hacknet/pb/smn_dict"
	"google.golang.org/protobuf/proto"
)

func callback(email string, serverAddr *net.UDPAddr, msg *hnp.ForwardMsg) (string, map[string]interface{}, error) {
	hnlog.Warn("get callback", hnlog.Fields{"email": email, "serverAddr": serverAddr, "msg": msg})

	return "", nil, nil
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func listen(port int) *net.UDPConn {
	res, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: port, Zone: ""})
	check(err)

	return res
}

func write(listener *net.UDPConn, port int, msg []byte) {
	_, err := listener.WriteToUDP(msg, &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: port, Zone: ""})
	check(err)
}

const (
	fromPort  = 501
	fromEmail = "007"
)

func main() {
	port := 500
	email := "root@suremoon.com"

	flag.IntVar(&port, "port", port, "local port")
	flag.StringVar(&email, "email", email, "server's email")
	flag.Parse()

	hacknetitf.NewServerForClient(port, email, callback)

	// the code for test.
	mfwd, err := hacknetitf.Pack_hnp_ForwardMsg(fromEmail, &hnp.ForwardMsg{
		FromEmail: fromEmail, FromPort: fromPort, Msg: "hello", FromIp: "127.0.0.5", Enums: int32(smn_dict.EDict_None),
	})

	check(err)
	mreg, err := hacknetitf.Pack_hnp_Register(fromEmail, &hnp.Register{PubKey: "??"})
	check(err)

	h1 := listen(fromPort)

	msg, err := proto.Marshal(mreg)
	check(err)
	write(h1, port, msg)
	msg, err = proto.Marshal(mfwd)
	check(err)
	write(h1, port, msg)

	for {
		time.Sleep(time.Second)
	}
}
