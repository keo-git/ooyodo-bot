package watcher

import (
	"io/ioutil"
	"log"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	gmail "google.golang.org/api/gmail/v1"

	"github.com/keo-git/ooyodo-bot/config"
	"github.com/keo-git/ooyodo-bot/utils"
)

type GmailWatcher struct {
	userId     string
	topic      string
	srv        *gmail.Service
	expiration int64
	historyId  uint64

	notifications *NotificationQueue
}

func NewGmailWatcher() *GmailWatcher {
	ctx := context.Background()
	conf := config.Config()
	secretFile, err := utils.AbsolutePath(conf.Credentials, conf.GmailSecret)
	if err != nil {
		log.Fatalf("Unable to get path to secret file: %v", err)
	}
	b, err := ioutil.ReadFile(secretFile)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client sercret file to config: %v", err)
	}

	client := getClient(ctx, config)
	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve gmail client: %v", err)
	}

	watcher := GmailWatcher{
		userId:        conf.UserId,
		topic:         conf.Topic,
		srv:           srv,
		expiration:    conf.Expiration,
		historyId:     conf.HistoryId,
		notifications: NewNotificationQueue(),
	}
	return &watcher
}

func (gw *GmailWatcher) StartWatcher() {
	now := time.Now().Unix() * 1000
	if now < gw.expiration {
		return
	}

	watchRequest := gmail.WatchRequest{TopicName: gw.topic}
	watchResponse, err := gw.srv.Users.Watch(gw.userId, &watchRequest).Do()
	if err != nil {
		log.Fatalf("Unable to set up watcher: %v", err)
	}
	gw.expiration = watchResponse.Expiration
	if gw.historyId == 0 {
		gw.historyId = watchResponse.HistoryId
	}
}

func (gw *GmailWatcher) Update() {
	now := time.Now().Unix() * 1000
	if now > gw.expiration {
		return
	}

	hisResp, err := gw.srv.Users.History.List(gw.userId).StartHistoryId(gw.historyId).Do()
	gw.historyId = hisResp.HistoryId
	if err != nil {
		log.Fatalf("Unable to retrieve history list: %v", err)
	}

	for _, history := range hisResp.History {
		for _, newMessage := range history.MessagesAdded {
			msg, _ := gw.srv.Users.Messages.Get(gw.userId, newMessage.Message.Id).Do()
			if isInbox(msg) {
				gw.notifications.Add(*msg)
			}
		}
	}
}

func isInbox(msg *gmail.Message) bool {
	for _, label := range msg.LabelIds {
		if label == "INBOX" {
			return true
		}
	}
	return false
}

func (gw *GmailWatcher) GetNotifications() []string {
	var notifications []string
	for !gw.notifications.IsEmpty() {
		notifications = append(notifications, gw.notifications.Get())
	}
	return notifications
}

func (gw *GmailWatcher) Close() {
	conf := config.Config()
	conf.Expiration = gw.expiration
	conf.HistoryId = gw.historyId
	conf.Close()
}
