package ooyodobot

import (
	"log"

	"github.com/keo-git/ooyodo-bot/config"
	"github.com/keo-git/ooyodo-bot/watcher"
	"gopkg.in/telegram-bot-api.v4"
)

type Ooyodo struct {
	*watcher.GmailWatcher
	api    *tgbotapi.BotAPI
	chatId int64
}

func NewOoyodo() *Ooyodo {
	conf := config.Config()
	api, err := tgbotapi.NewBotAPI(conf.TelegramToken)
	if err != nil {
		log.Fatalf("Unable to initialize bot API: %v", err)
	}
	log.Printf("Authorized on account %s\n", api.Self.UserName)

	return &Ooyodo{watcher.NewGmailWatcher(), api, conf.ChatId}
}

func (ooyodo *Ooyodo) SendNotification(notification string) {
	telMsg := tgbotapi.NewMessage(ooyodo.chatId, notification)
	ooyodo.api.Send(telMsg)
}
