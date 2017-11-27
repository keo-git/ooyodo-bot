package watcher

import (
	"time"

	"github.com/keo-git/go-bot/handler"
	gmail "google.golang.org/api/gmail/v1"
)

type GmailWatcher struct {
	done chan struct{}

	userId string
	srv    *gmail.Service

	sub *subscription
}

func NewGmailWatcher(secret, token, sub, userId string) (*GmailWatcher, error) {
	client, err := getClient(secret, token)
	if err != nil {
		return nil, err
	}
	srv, err := gmail.New(client)
	if err != nil {
		return nil, err
	}

	s, err := NewSubscription(srv, userId, sub)
	if err != nil {
		return nil, err
	}

	return &GmailWatcher{
		done:   make(chan struct{}),
		userId: userId,
		srv:    srv,
		sub:    s,
	}, nil
}

func (gw *GmailWatcher) Start(updates chan<- *handler.Message, errc chan<- error) {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			if gw.sub.IsExpired() {
				if err := gw.sub.Subscribe(gw.srv, gw.userId); err != nil {
					errc <- err
					return
				}
			}

			hisResp, err := gw.srv.Users.History.List(gw.userId).StartHistoryId(gw.sub.HistoryId).Do()
			if err != nil {
				errc <- err
				return
			}

			for _, history := range hisResp.History {
				for _, newMessage := range history.MessagesAdded {
					if isInbox(newMessage.Message) {
						msg, err := gw.srv.Users.Messages.Get(gw.userId, newMessage.Message.Id).Do()
						if err != nil {
							errc <- err
							continue
						}

						headers := getMessageHeaders(msg, "Date", "From", "To", "Subject")
						body := getMessageBodyText(msg)
						attachments := getMessageAttachments(msg, gw.srv, gw.userId, msg.Id)

						n := NewNotification(headers, body, attachments)
						updates <- n

						if msg.HistoryId > gw.sub.HistoryId {
							gw.sub.HistoryId = msg.HistoryId
							gw.sub.Save()
						}
					}
				}
			}
		case <-gw.done:
			ticker.Stop()
		}
	}
}

func (gw GmailWatcher) Stop() {
	gw.done <- struct{}{}
}
