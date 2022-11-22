package main

import (
	"flag"
	"fmt"
	"net"
	"time"

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

	retry := true
	client.RegisterMsgHandler("NoticeS", func(addr *net.UDPAddr, params gkritf.GeekerMsg) (err error) {
		fmt.Println("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$", params)
		var noticeS gkritf.NoticeS
		if noticeS, err = gkritf.NewNoticeS(params); err != nil {
			return err
		}

		conn := client.GetConn()
		conn.WriteToUDP([]byte(gkritf.Hole{Msg: "hole ~"}.Message()),
			noticeS.NodeInfo.BuildUDPAddr())
		return nil
	})

	client.RegisterMsgHandler("Hole", func(addr *net.UDPAddr, params gkritf.GeekerMsg) (err error) {
		retry = false
		fmt.Println("hole success, another addr = ", addr)
		return nil
	})

	client.Connect(*localPort, *remoteIP, *remotePort)
	for retry {
		client.Send(gkritf.NoticeC{
			TargetSession: *targetSession,
			ExtraData:     "",
		}.Message())
		time.Sleep(time.Second * 5)
	}
	ch <- 1
}
