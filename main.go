package main

import (
	"flag"
	"log"
	"path/filepath"

	"github.com/keo-git/ooyodo-bot/config"
	"github.com/keo-git/ooyodo-bot/ooyodobot"
	"github.com/keo-git/ooyodo-bot/utils"
)

var file = flag.String("config", "", "path to config file")

func main() {
	log.Println("Starting Ooyodo routine...")
	flag.Parse()
	abs, err := filepath.Abs(*file)
	if err != nil {
		log.Fatalf("Unable to open config file: %v", err)
	}
	conf, err := config.Config(abs)
	if err != nil {
		log.Fatalf("Unable to create config instance: %v", err)
	}
	ooyodo, err := ooyodobot.NewOoyodo(conf.GmailSecret, conf.GmailToken,
		conf.Subscription, conf.UserId, conf.TelegramToken, conf.ChatId)
	if err != nil {
		log.Fatalf("Unable to create Ooyodo instance: %v", err)
	}

	c := make(chan interface{}, 1)
	go utils.Log(c)
	go utils.Ping()

	for {
		err = ooyodo.Update()
		if err != nil {
			c <- err
			continue
		}
		notifications := ooyodo.GetNotifications()
		if len(notifications) > 0 {
			c <- len(notifications)
		}
		for _, n := range notifications {
			ooyodo.Notify(*n)
		}
	}
}
