package main

import (
	"flag"

	"github.com/ProtossGenius/hacknet/gkritf"
)

func main() {
	port := flag.Int("p", 998, "local port")
	flag.Parse()
	server := gkritf.NewGeekerNetUDPServer()
	server.Listen(*port)

	ch := make(chan int, 0)
	ch <- 1
}
