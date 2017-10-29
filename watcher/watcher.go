package watcher

import (
	"io/ioutil"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	gmail "google.golang.org/api/gmail/v1"
	//"github.com/keo-git/ooyodo-bot/config"
)

type GmailWatcher struct {
	userId string
	srv    *gmail.Service

	sub *subscription
	nq  *NotificationQueue
}

func NewGmailWatcher(secret, token, sub, userId string) (*GmailWatcher, error) {
	ctx := context.Background()

	b, err := ioutil.ReadFile(secret)
	if err != nil {
		return nil, err
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		return nil, err
	}

	client := getClient(ctx, config, token)
	srv, err := gmail.New(client)
	if err != nil {
		return nil, err
	}

	return &GmailWatcher{
		userId: userId,
		srv:    srv,
		sub:    NewSubscription(srv, userId, sub),
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
			msg, err := gw.srv.Users.Messages.Get(gw.userId, newMessage.Message.Id).Do()
			if err != nil {
				return err
			}
			if isInbox(msg) {
				if msg.HistoryId > gw.sub.HistoryId {
					gw.sub.HistoryId = msg.HistoryId
					gw.sub.Save()
				}
				gw.nq.Push(*msg)
			}
		}
	}
	//gw.sub.HistoryId = hisResp.HistoryId
	//gw.sub.Save()
	return nil
}

func (gw *GmailWatcher) GetNotification() *notification {
	return gw.nq.Pop()
}

func (gw *GmailWatcher) GetNotifications() []*notification {
	var ns []*notification
	for !gw.nq.IsEmpty() {
		ns = append(ns, gw.GetNotification())
	}
	return ns
}

func isInbox(msg *gmail.Message) bool {
	for _, label := range msg.LabelIds {
		if label == "INBOX" {
			return true
		}
	}
	return false
}

/*func (gw *GmailWatcher) Close() {
	conf := config.Config()
	conf.Expiration = gw.expiration
	conf.HistoryId = gw.historyId
	conf.Close()
}*/
