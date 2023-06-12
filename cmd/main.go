package main

import (
	"chat-app/cmd/server"
	"chat-app/config"
	"github.com/sirupsen/logrus"
)

func initLogger() {
	logrus.SetFormatter(
		&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		},
	)

	if logLvl, err := logrus.ParseLevel(config.GetAppConfig().Logger.Level); err != nil {
		logrus.SetLevel(logLvl)
	}
}

func main() {
	config.Load()
	initLogger()

	httpServer := server.HTTPServer{}
	httpServer.Init()
	if err := httpServer.Run(); err != nil {
		logrus.WithError(err).Fatal("failed to start httpServer")
	}
}
