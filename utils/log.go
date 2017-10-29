package utils

import (
	"log"
	"time"
)

func Log(c chan interface{}) {
	for {
		p := <-c
		switch p := p.(type) {
		case error:
			log.Printf("Error occured: %v", p)
		case string:
			log.Println(p)
		case int:
			log.Printf("%d new message(s).\nSending notification(s)...", p)
		default:
			log.Printf("%v", p)
		}
	}
}

func Ping() {
	for {
		time.Sleep(time.Minute * 15)
		log.Printf("Daijoubu!")
	}
}
