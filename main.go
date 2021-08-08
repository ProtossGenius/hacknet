package main

import (
	"github.com/ProtossGenius/hacknet/hnlog"
	"github.com/sirupsen/logrus"
)

func main() {
	hnlog.Warn("haha", logrus.Fields{})
}
