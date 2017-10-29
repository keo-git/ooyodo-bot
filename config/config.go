package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type config struct {
	UserId       string `json:"gmail_id"`
	GmailSecret  string `json:"gmail_secret"`
	GmailToken   string `json:"gmail_token"`
	Subscription string `json:"gmail_subscription"`

	TelegramToken string `json:"telegram_token"`
	ChatId        int64  `json:"chat_id"`
}

var confInstance *config = nil

func Config(file string) (*config, error) {
	if confInstance == nil {
		confInstance = new(config)
		f, err := os.Open(file)
		if err != nil {
			confInstance = nil
			return nil, err
		}
		err = json.NewDecoder(f).Decode(confInstance)
		if err != nil {
			confInstance = nil
			return nil, err
		}
		confInstance.GmailSecret, err = filepath.Abs(confInstance.GmailSecret)
		if err != nil {
			confInstance = nil
			return nil, err
		}
		confInstance.GmailToken, err = filepath.Abs(confInstance.GmailToken)
		if err != nil {
			confInstance = nil
			return nil, err
		}
		confInstance.Subscription, err = filepath.Abs(confInstance.Subscription)
		if err != nil {
			confInstance = nil
			return nil, err
		}
	}
	return confInstance, nil
}
