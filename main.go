package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/keo-git/ooyodo-bot/config"
	"github.com/keo-git/ooyodo-bot/ooyodobot"
)

func sigHandler(c chan os.Signal, ooyodo *ooyodobot.Ooyodo) {
	switch <-c {
	case syscall.SIGTERM:
		ooyodo.Close()
		os.Exit(0)
	case syscall.SIGINT:
		ooyodo.Close()
		os.Exit(0)
	}
}

var file = flag.String("config", "ooyodo-bot.json", "path to config file")

func init() {

}

func main() {
	flag.Parse()
	config.InitConfig(*file)
	ooyodo := ooyodobot.NewOoyodo()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	go sigHandler(c, ooyodo)

	ooyodo.StartWatcher()
	defer ooyodo.Close()

	for {
		ooyodo.Update()
		notifications := ooyodo.GetNotifications()
		for _, notification := range notifications {
			ooyodo.SendNotification(notification)
		}
	}
}
