package config

import (
	"encoding/json"
	"log"
	"os"
)

type config struct {
	File string

	UserId string `json:"user_id"`
	Topic  string `json:"topic"`

	Credentials string `json:"cdir"`
	GmailSecret string `json:"gmail_secret"`
	GmailToken  string `json:"gmail_token"`

	TelegramToken string `json:"telegram_token"`
	ChatId        int64  `json:"chat_id"`

	Expiration int64  `json:"exp_date"`
	HistoryId  uint64 `json:"history_id"`
}

var confInstance *config = nil

func InitConfig(file string) {
	confInstance = new(config)
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("Unable to open config file %v", err)
	}
	err = json.NewDecoder(f).Decode(confInstance)
	if err != nil {
		log.Fatalf("Unable to parse config file %v", err)
	}
	confInstance.File = file
}

func Config() *config {
	return confInstance
}

func (c *config) Close() {
	f, err := os.OpenFile(c.File, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to open config file %v", err)
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(&c)
	if err != nil {
		log.Fatalf("Unable to write to config file %v", err)
	}
}
