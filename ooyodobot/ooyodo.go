package ooyodobot

import (
	"github.com/keo-git/go-bot/bot"
	"github.com/keo-git/go-bot/handler"
	"github.com/keo-git/go-bot/handler/telegram"
	"github.com/keo-git/ooyodo-bot/watcher"
)

type Ooyodo struct {
	*bot.Bot
	*watcher.GmailWatcher
	updates chan *handler.Message
	errc    chan error
	done    chan struct{}
	chatId  int64
}

func NewOoyodo(gmailSecret, gmailToken, sub, userId, telToken string, chatId int64) (*Ooyodo, error) {
	tel, err := telegram.NewTelegramHandler(telToken)
	if err != nil {
		return nil, err
	}

	w, err := watcher.NewGmailWatcher(gmailSecret, gmailToken, sub, userId)
	if err != nil {
		return nil, err
	}

	return &Ooyodo{
		Bot:          bot.NewBot(tel),
		GmailWatcher: w,
		updates:      make(chan *handler.Message),
		errc:         make(chan error),
		done:         make(chan struct{}),
		chatId:       chatId,
	}, nil
}

func (o Ooyodo) Start() {
	go o.Bot.Start()
	go o.GmailWatcher.Start(o.updates, o.errc)
	go o.notify()
}

func (o Ooyodo) Stop() {
	o.done <- struct{}{}
}

func (o Ooyodo) notify() {
	for {
		select {
		case update := <-o.updates:
			update.Channel.Id = o.chatId
			go o.Send(update)
		case <-o.done:
			o.Bot.Close()
			o.GmailWatcher.Stop()
			close(o.errc)
			return
			//case err := <-o.errc:

		}
	}
}
