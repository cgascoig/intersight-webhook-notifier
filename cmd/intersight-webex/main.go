package main

import (
	"github.com/cgascoig/intersight-webex/pkg/iswbx"
	"github.com/sirupsen/logrus"
)

func main() {
	server := iswbx.NewServer()
	logrus.SetLevel(logrus.DebugLevel)

	err := server.Run()
	if err != nil {
		logrus.Fatal(err)
	}
}
