package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/keo-git/ooyodo-bot/config"
	"github.com/keo-git/ooyodo-bot/ooyodobot"
)

func sigHandler(c chan os.Signal, ooyodo *ooyodobot.Ooyodo) {
	switch <-c {
	case syscall.SIGTERM:
		//ooyodo.Close()
		os.Exit(0)
	case syscall.SIGINT:
		//ooyodo.Close()
		os.Exit(0)
	}
}

var file = flag.String("config", "", "path to config file")

func init() {

}

func main() {
	flag.Parse()
	abs, err := filepath.Abs(*file)
	if err != nil {
		log.Fatalf("Unable to open config file: %v", err)
	}
	conf, err := config.Config(abs)
	if err != nil {
		log.Fatalf("Unable to open create config instance: %v", err)
	}
	ooyodo, err := ooyodobot.NewOoyodo(conf.GmailSecret, conf.GmailToken,
		conf.Subscription, conf.UserId, conf.TelegramToken, conf.ChatId)
	if err != nil {
		log.Fatalf("Unable to open create Ooyodo instance: %v", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	go sigHandler(c, ooyodo)

	//ooyodo.StartWatcher()
	//defer ooyodo.Close()

	for {
		err = ooyodo.Update()
		if err != nil {
			log.Printf("Unable to update: %v", err)
			continue
		}
		notifications := ooyodo.GetNotifications()
		for _, n := range notifications {
			fmt.Println(n.MsgText)
			ooyodo.Notify(*n)
		}
	}
}
