package main

import (
	"os"

	"github.com/cgascoig/intersight-webhook-notifier/pkg/iswbx"
	"github.com/sirupsen/logrus"
)

func main() {
	server := iswbx.NewServer()

	logrus.SetFormatter(&logrus.JSONFormatter{})

	if os.Getenv("DEBUG") != "" {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("Logging set to DEBUG")
	}

	err := server.Run()
	if err != nil {
		logrus.Fatal(err)
	}
}
