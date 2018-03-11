package watcher

import (
	"github.com/svarw/ooyodo-bot/bot/handler"
	"gopkg.in/telegram-bot-api.v4"
)

func NewNotification(headers map[string]string, body string, attachments []tgbotapi.FileBytes) *bot.Message {
	date := headers["Date"]
	from := headers["From"]
	to := headers["To"]
	subject := headers["Subject"]
	msg := multiString(date, "From: "+from, "To: "+to, "Subject: "+subject, body)
	return &bot.Message{
		Msg:   msg,
		Files: attachments,
		//Channel: ,
	}
}

func multiString(args ...string) string {
	var ms string
	for _, s := range args {
		ms += s + "\n"
	}
	return ms
}
