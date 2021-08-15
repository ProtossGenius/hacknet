package main

import (
	"net"

	"github.com/ProtossGenius/hacknet/hacknetitf"
	"github.com/ProtossGenius/hacknet/hnlog"
	"github.com/ProtossGenius/hacknet/pb/hnp"
)

func callback(email string, serverAddr *net.UDPAddr, msg *hnp.ForwardMsg) (string, map[string]interface{}, error) {
	hnlog.Info("get callback", hnlog.Fields{"": ""})

	return "", nil, nil
}

func main() {
	hacknetitf.NewServerForClient(500, "root@suremoon.com", callback)
}
