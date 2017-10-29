package ooyodobot

import (
	"log"

	"github.com/keo-git/ooyodo-bot/watcher"
	"gopkg.in/telegram-bot-api.v4"
)

type Ooyodo struct {
	*watcher.GmailWatcher
	api    *tgbotapi.BotAPI
	chatId int64
}

func NewOoyodo(gmailSecret, gmailToken, sub, userId, telToken string, chatId int64) (*Ooyodo, error) {
	api, err := tgbotapi.NewBotAPI(telToken)
	if err != nil {
		return nil, err
	}
	log.Printf("Authorized on account %s\n", api.Self.UserName)

	w, err := watcher.NewGmailWatcher(gmailSecret, gmailToken, sub, userId)
	return &Ooyodo{w, api, chatId}, err
}

func (ooyodo *Ooyodo) Notify(n watcher.Notification) {
	telMsg := tgbotapi.NewMessage(ooyodo.chatId, n.MsgText)
	ooyodo.api.Send(telMsg)
	for _, msgFile := range n.MsgFiles {
		docUpload := tgbotapi.NewDocumentUpload(ooyodo.chatId, msgFile)
		ooyodo.api.Send(docUpload)
	}
}
