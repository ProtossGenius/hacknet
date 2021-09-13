package main

import (
	"flag"
	"fmt"
	"net"
)

// createConn create Conn.
func createConn(localPort int, remoteIP string, remotePort int) (*net.UDPConn, error) {
	local := &net.UDPAddr{IP: net.IPv4zero, Port: localPort, Zone: ""}
	if remoteIP != "" { // bind and conn
		return net.DialUDP("udp",
			local,
			&net.UDPAddr{IP: net.ParseIP(remoteIP), Port: remotePort, Zone: ""})
	}

	// only bind
	return net.ListenUDP("udp", local)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// listen listen msg and print.
func listen(conn *net.UDPConn) {
	data := make([]byte, 4096)

	for {
		_, addr, err := conn.ReadFromUDP(data)
		fmt.Println(addr, string(data), err)
	}
}

func main() {
	localPort := flag.Int("localPort", 2001, "local udp port to listen, -1 not listen")
	remotePort := flag.Int("remotePort", 2002, "remote port")
	remoteIP := flag.String("remoteIP", "", "remote IP")

	flag.Parse()

	conn, err := createConn(*localPort, *remoteIP, *remotePort)
	check(err)

	go listen(conn)

	line := ""

	for {
		fmt.Scanln(&line)

		if _, err := conn.Write([]byte(line)); err != nil {
			fmt.Println("err = ", err)
		}
	}
}
