package watcher

import (
	gmail "google.golang.org/api/gmail/v1"
)

//msgText in format
//Date
//From: sender@email1.com
//To : receiver@email2.com
//Subject
//Body
type notification struct {
	msgText string
	//msgFiles ???
}

func NewNotification(msg gmail.Message) *notification {
	var date, from, to, subject, body string
	for _, header := range msg.Payload.Headers {
		switch header.Name {
		case "Date":
			date = header.Value
		case "From":
			from = header.Value
		case "To":
			to = header.Value
		case "Subject":
			subject = header.Value
		}
	}
	msgText := multiString(date, "From: "+from, "To: "+to, "Subject: "+subject, body)
	return &notification{msgText}
}

func (n notification) GetMsgText() string {
	return n.msgText
}

func multiString(args ...string) string {
	var ms string
	for _, s := range args {
		ms += s + "\n"
	}
	return ms
}

type NotificationQueue struct {
	notifications []*notification
}

func NewNotificationQueue() *NotificationQueue {
	return &NotificationQueue{}
}

func (nq *NotificationQueue) Push(msg gmail.Message) {
	n := NewNotification(msg)
	nq.notifications = append(nq.notifications, n)
}

func (nq *NotificationQueue) Pop() *notification {
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
