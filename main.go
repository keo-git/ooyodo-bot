package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/keo-git/ooyodo-bot/config"
	"github.com/keo-git/ooyodo-bot/ooyodobot"
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

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ooyodo.Start()
	for range c {
		ooyodo.Stop()
		close(c)
	}
}
