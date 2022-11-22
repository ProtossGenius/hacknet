package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/ProtossGenius/hacknet/gkritf"
)

func main() {
	localPort := flag.Int("lp", 999, "local port")
	remoteIP := flag.String("r", "127.0.0.1", "remote ip")
	remotePort := flag.Int("rp", 998, "remote port")
	session := flag.String("ssid", "0xCF", "target session id")
	targetSession := flag.String("tssid", "0xCF", "target session id")
	flag.Parse()
	ch := make(chan int, 0)
	client := gkritf.NewGeekerNetUDPClient(*session)
	client.RegisterMsgHandler("Notice", func(addr *net.UDPAddr, params string) error {
		if params == "done" {
			return nil
		}

		fmt.Println(params)

		return nil
	})
	client.Connect(*localPort, *remoteIP, *remotePort)
	client.Send("Notice#" + *targetSession)
	ch <- 1
}
