package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	UserID       string `json:"gmail_id"`
	GmailSecret  string `json:"gmail_secret"`
	GmailToken   string `json:"gmail_token"`
	Subscription string `json:"gmail_subscription"`

	TelegramToken string `json:"telegram_token"`
	ChatID        int64  `json:"chat_id"`
}

func NewConfig(file string) (config *Config, err error) {
	config = new(Config)
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(f).Decode(config)
	if err != nil {
		return nil, err
	}
	config.GmailSecret, err = filepath.Abs(config.GmailSecret)
	if err != nil {
		return nil, err
	}
	config.GmailToken, err = filepath.Abs(config.GmailToken)
	if err != nil {
		return nil, err
	}
	config.Subscription, err = filepath.Abs(config.Subscription)
	if err != nil {
		return nil, err
	}
	return
}
