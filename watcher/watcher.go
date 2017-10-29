package watcher

import (
	gmail "google.golang.org/api/gmail/v1"
)

type GmailWatcher struct {
	userId string
	srv    *gmail.Service

	sub *subscription
	nq  *NotificationQueue
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
		userId: userId,
		srv:    srv,
		sub:    s,
		nq:     NewNotificationQueue(),
	}, nil
}

func (gw *GmailWatcher) Update() error {
	if gw.sub.IsExpired() {
		err := gw.sub.Subscribe(gw.srv, gw.userId)
		if err != nil {
			return err
		}
	}

	hisResp, err := gw.srv.Users.History.List(gw.userId).StartHistoryId(gw.sub.HistoryId).Do()
	if err != nil {
		return err
	}

	for _, history := range hisResp.History {
		for _, newMessage := range history.MessagesAdded {
			if isInbox(newMessage.Message) {
				msg, err := gw.srv.Users.Messages.Get(gw.userId, newMessage.Message.Id).Do()
				if err != nil {
					return err
				}

				headers := getMessageHeaders(msg, "Date", "From", "To", "Subject")
				body := getMessageBodyText(msg)
				attachments := getMessageAttachments(msg, gw.srv, gw.userId, msg.Id)

				n := NewNotification(headers, body, attachments)
				gw.nq.Push(n)

				if msg.HistoryId > gw.sub.HistoryId {
					gw.sub.HistoryId = msg.HistoryId
					gw.sub.Save()
				}
			}
		}
	}
	return nil
}

func (gw *GmailWatcher) GetNotification() *Notification {
	return gw.nq.Pop()
}

func (gw *GmailWatcher) GetNotifications() []*Notification {
	var ns []*Notification
	for !gw.nq.IsEmpty() {
		ns = append(ns, gw.GetNotification())
	}
	return ns
}
