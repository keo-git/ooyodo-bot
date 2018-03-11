package bot

import (
	"gopkg.in/telegram-bot-api.v4"
)

type Message struct {
	Msg     string
	Files   []tgbotapi.FileBytes
	Channel int64
}

type Handler struct {
	api  *tgbotapi.BotAPI
	done chan struct{}
}

func NewHandler(token string) (*Handler, error) {
	api, err := tgbotapi.NewBotAPI(token)
	return &Handler{api: api, done: make(chan struct{})}, err
}

func (h Handler) Send(msg Message, errc chan<- error) {
	msgConfig := tgbotapi.NewMessage(msg.Channel, msg.Msg)
	_, err := h.api.Send(msgConfig)
	if err != nil {
		errc <- err
		return
	}
	for _, file := range msg.Files {
		docUpload := tgbotapi.NewDocumentUpload(msg.Channel, file)
		_, err := h.api.Send(docUpload)
		if err != nil {
			errc <- err
			return
		}
	}
}
