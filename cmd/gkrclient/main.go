package main

import (
	"flag"
	"fmt"
	"log"
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
	hole := flag.Bool("hole", false, "start hole.")
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

		client.Send(gkritf.NoticeC{TargetSession: *targetSession, ExtraData: noticeS.ExtraData + "1"}.Message())

		if len(noticeS.ExtraData) > 0 {
			retry = false
		}

		if len(noticeS.ExtraData) > 1 {
			conn := client.MoveConn()
			conn.Close()
			go doHole(*localPort, noticeS.NodeInfo.BuildUDPAddr())
		}

		return nil
	})

	client.Connect(*localPort, *remoteIP, *remotePort)
	for retry && *hole {
		client.Send(gkritf.NoticeC{
			TargetSession: *targetSession,
			ExtraData:     "",
		}.Message())
		time.Sleep(time.Second * 5)
	}
	ch <- 1
}
func check(err error) {
	if err != nil {
		panic(err)
	}
}

// hole do hole .
func doHole(localPort int, targetAddr *net.UDPAddr) {
	conn, err := net.DialUDP("udp", &net.UDPAddr{IP: net.ParseIP("8.8.8.8"), Port: localPort}, targetAddr)
	if err != nil {
		panic(err)
	}

	go func() {
		data := make([]byte, 1024)
		for {
			n, addr, err := conn.ReadFromUDP(data)
			check(err)
			if n > 0 {
				log.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@2 get msg from ", addr, ":        ", string(data[:n]))
			}
		}
	}()

	for {
		log.Println("send hole message from 0:", localPort, " to ", targetAddr, "message is ", "hole ~")
		_, err := conn.Write([]byte(gkritf.Hole{Msg: "hole ~"}.Message()))
		check(err)
		time.Sleep(time.Second)
	}
}
