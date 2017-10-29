package watcher

import (
	"gopkg.in/telegram-bot-api.v4"
)

//msgText in format
//Date
//From: sender@email1.com
//To : receiver@email2.com
//Subject: subject
//Body
type Notification struct {
	MsgText  string
	MsgFiles []tgbotapi.FileBytes
}

func NewNotification(headers map[string]string, body string, attachments []tgbotapi.FileBytes) *Notification {
	date := headers["Date"]
	from := headers["From"]
	to := headers["To"]
	subject := headers["Subject"]
	msgText := multiString(date, "From: "+from, "To: "+to, "Subject: "+subject, body)
	return &Notification{msgText, attachments}
}

func multiString(args ...string) string {
	var ms string
	for _, s := range args {
		ms += s + "\n"
	}
	return ms
}

type NotificationQueue struct {
	notifications []*Notification
}

func NewNotificationQueue() *NotificationQueue {
	return &NotificationQueue{}
}

func (nq *NotificationQueue) Push(n *Notification) {
	nq.notifications = append(nq.notifications, n)
}

func (nq *NotificationQueue) Pop() *Notification {
	if nq.IsEmpty() {
		return nil
	}
	n := nq.notifications[0]
	nq.notifications = nq.notifications[1:]
	return n
}

func (n NotificationQueue) IsEmpty() bool {
	if len(n.notifications) > 0 {
		return false
	}
	return true
}
