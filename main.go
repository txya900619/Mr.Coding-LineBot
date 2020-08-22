package main

import (
	"Mr.Coding-LineBot/config"
	"Mr.Coding-LineBot/spreadsheets"
	"github.com/line/line-bot-sdk-go/linebot"
	"google.golang.org/api/sheets/v4"
	"log"
	"net/http"
)

var bot *linebot.Client
var googleService *sheets.Service

func main() {
	// Get config at ./config.yaml
	c, err := config.New()
	if err != nil {
		log.Fatalf("Read config.yaml file fail, %v", err)
	}

	googleService, err = spreadsheets.New()
	if err != nil {
		log.Fatal(err)
	}

	bot, err = linebot.New(c.ChannelSecret, c.ChannelToken)
	if err != nil {
		log.Fatalf("Create linebot fail, %v", err)
	}

	http.HandleFunc("/callback", callbackHandler)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		switch event.Type {
		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if message.Text == "Mr.Coding 表單" {
					return
				}
			}
		}
	}
}
