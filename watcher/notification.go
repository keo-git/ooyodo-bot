package watcher

import (
	gmail "google.golang.org/api/gmail/v1"
)

func newMultiString(args ...string) string {
	var ms string
	for _, s := range args {
		ms += s + "\n"
	}
	return ms
}

type NotificationQueue struct {
	notifications []string
}

func NewNotificationQueue() *NotificationQueue {
	return &NotificationQueue{}
}

func (n *NotificationQueue) Add(msg gmail.Message) {

	var date, from string
	for _, header := range msg.Payload.Headers {
		switch header.Name {
		case "Date":
			date = header.Value
		case "From":
			from = header.Value
		}
	}
	newNotif := newMultiString(date, from, msg.Snippet)
	n.notifications = append(n.notifications, newNotif)
}

func (n *NotificationQueue) Get() string {
	if n.IsEmpty() {
		return ""
	}
	message := n.notifications[0]
	n.notifications = n.notifications[1:]
	return message
}

func (n NotificationQueue) IsEmpty() bool {
	if len(n.notifications) > 0 {
		return false
	}
	return true
}
