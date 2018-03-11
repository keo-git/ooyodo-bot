package ooyodobot

import (
	"github.com/svarw/ooyodo-bot/bot/handler"
	"github.com/svarw/ooyodo-bot/config"
	"github.com/svarw/ooyodo-bot/watcher"
)

type Ooyodo struct {
	*watcher.GmailWatcher
	*bot.Handler
	updates chan *bot.Message
	errc    chan error
	done    chan struct{}
	chatID  int64
}

func NewOoyodo(c config.Config) (*Ooyodo, error) {
	h, err := bot.NewHandler(c.TelegramToken)
	if err != nil {
		return nil, err
	}

	w, err := watcher.NewGmailWatcher(c.GmailSecret, c.GmailToken, c.Subscription, c.UserID)
	if err != nil {
		return nil, err
	}

	return &Ooyodo{
		//Bot:          bot.NewBot(tel),
		GmailWatcher: w,
		Handler:      h,
		updates:      make(chan *bot.Message),
		errc:         make(chan error),
		done:         make(chan struct{}),
		chatID:       c.ChatID,
	}, nil
}

func (o Ooyodo) Start() {
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
			update.Channel = o.chatID
			go o.Send(*update, o.errc)
		case <-o.done:
			o.GmailWatcher.Stop()
			close(o.errc)
			return
			//case err := <-o.errc:
		}
	}
}
