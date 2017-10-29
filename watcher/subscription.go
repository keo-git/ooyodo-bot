package watcher

import (
	"encoding/json"
	"os"

	"github.com/keo-git/ooyodo-bot/utils"
	gmail "google.golang.org/api/gmail/v1"
)

type subscription struct {
	file       string
	Topic      string `json:"topic"`
	Expiration int64  `json:"expiration"`
	HistoryId  uint64 `json:"history_id"`
}

func NewSubscription(srv *gmail.Service, userId, file string) (*subscription, error) {
	s, err := subscriptionFromFile(file)
	s.file = file
	if err != nil || s.IsExpired() {
		if err = s.Subscribe(srv, userId); err != nil {
			return nil, err
		}
		s.Save()
	}
	return s, nil
}

func (s *subscription) Subscribe(srv *gmail.Service, userId string) error {
	if !s.IsExpired() {
		return nil
	}
	watchRequest := gmail.WatchRequest{TopicName: s.Topic}
	watchResponse, err := srv.Users.Watch(userId, &watchRequest).Do()
	if err != nil {
		return err
	}
	s.Expiration = watchResponse.Expiration
	s.HistoryId = watchResponse.HistoryId
	return nil
}

func (s subscription) IsExpired() bool {
	now := utils.UnixMili()
	if now > s.Expiration {
		return true
	}
	return false
}

func subscriptionFromFile(file string) (*subscription, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	s := &subscription{}
	err = json.NewDecoder(f).Decode(s)
	return s, err
}

func (s *subscription) Save() error {
	f, err := os.OpenFile(s.file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(s)
	return err
}
