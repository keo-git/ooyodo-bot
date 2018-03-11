package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/svarw/ooyodo-bot/bot"
	"github.com/svarw/ooyodo-bot/config"
)

var file = flag.String("config", "", "path to config file")

func main() {
	log.Println("Starting Ooyodo...")
	flag.Parse()
	abs, err := filepath.Abs(*file)
	if err != nil {
		log.Fatalf("Unable to open config file: %v", err)
	}
	conf, err := config.NewConfig(abs)
	if err != nil {
		log.Fatalf("Unable to create config instance: %v", err)
	}
	ooyodo, err := ooyodobot.NewOoyodo(*conf)
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
