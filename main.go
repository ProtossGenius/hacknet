package main

import (
	"flag"
	"net"
	"time"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_file"
	"github.com/ProtossGenius/hacknet/hacknetitf"
	"github.com/ProtossGenius/hacknet/hnlog"
	"github.com/ProtossGenius/hacknet/pb/hnep"
	"github.com/ProtossGenius/hacknet/pb/hnp"
	"github.com/ProtossGenius/hacknet/pb/smn_dict"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func NotOnForwardMsg(email string, serverAddr *net.UDPAddr, msg *hnp.ForwardMsg) (
	string, map[string]interface{}, error) {
	hnlog.Info("get callback", hnlog.Fields{"email": email, "serverAddr": serverAddr, "msg": msg})

	return "", nil, nil
}

func ClientOnForwardMsg(email string, serverAddr *net.UDPAddr, msg *hnp.ForwardMsg) (
	string, map[string]interface{}, error) {
	hnlog.Info("get callback", hnlog.Fields{"email": email, "serverAddr": serverAddr, "msg": msg})
	_, detail, err := server.Register(email, serverAddr, &hnp.Register{})
	if err != nil {
		detail["err"] = err
		hnlog.Error("ClientOnForwardMsg#server.Register", detail)
	}

	return "", nil, nil
}

func NotOnResult(email string, serverAddr *net.UDPAddr, msg *hnp.Result) (
	string, map[string]interface{}, error) {
	hnlog.Info("get callback", hnlog.Fields{"email": email, "serverAddr": serverAddr, "msg": msg})

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

var (
	targetIP   = "127.0.0.1"
	targetPort = 501

	server hacknetitf.ServerItf
)

func main() {
	port := 500
	email := "root@suremoon.com"
	pubKeyPath := ""

	flag.IntVar(&port, "hnport", port, "hack net's local port")
	flag.StringVar(&email, "hnemail", email, "hack net's local server's email")
	flag.StringVar(&targetIP, "hntip", targetIP, "hack net's target ip")
	flag.IntVar(&targetPort, "hntport", port, "hack net's target port")
	flag.StringVar(&pubKeyPath, "pubkey", pubKeyPath, "pubkey file's path")
	flag.Parse()

	pubKey := ""

	if pubKeyPath != "" {
		data, err := smn_file.FileReadAll(pubKeyPath)
		check(err)

		pubKey = string(data)
	}

	server = hacknetitf.NewServer(port, email, pubKey, NotOnForwardMsg, NotOnResult)

	anyMsg, err := anypb.New(&hnep.StrMsg{Msg: "hello"})
	check(err)

	// the code for test.
	mfwd, err := hacknetitf.Pack_hnp_ForwardMsg(fromEmail, &hnp.ForwardMsg{
		FromEmail: fromEmail, FromPort: fromPort, Msg: anyMsg, FromIp: "127.0.0.5", Enums: int32(smn_dict.EDict_None),
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
