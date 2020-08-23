package main

import (
	"Mr.Coding-LineBot/config"
	"Mr.Coding-LineBot/mrcoding"
	"fmt"
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
	"net/http"
)

var bot *mrcoding.Bot

func main() {
	// Get config at ./config.yaml
	c, err := config.New()
	if err != nil {
		log.Fatalf("Read config.yaml file fail, %v", err)
	}

	bot, err = mrcoding.New(c)
	if err != nil {
		log.Fatalf("Create linebot fail, %v", err)
	}
	http.HandleFunc("/callback", callbackHandler)
	err = http.ListenAndServe(":1225", nil)
	fmt.Println("serve on :1225")
	if err != nil {
		log.Fatal(err)
	}
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
				switch message.Text {
				case "Mr.Coding 表單":
					message, err := bot.QuestionStart(event.Source.UserID)
					if err != nil {
						log.Fatal(err)
					}

					bot.ReplyMessage(event.ReplyToken, message).Do()
					return
				case "社團博覽會有獎徵答":
					return
				default:
					answerRowID, err := bot.FindAnswerRowID(event.Source.UserID)
					if err != nil {
						log.Fatal(err)
					}

					if answerRowID != 0 {
						message, err := bot.SaveAnswerAndGetNextMessage(message.Text, answerRowID, event.Source.UserID)
						if err != nil {
							log.Fatal(err)
						}

						bot.ReplyMessage(event.ReplyToken, message).Do()
						return
					}
				}
			}
		case linebot.EventTypePostback:
			switch event.Postback.Data {
			case "pass":
				answerRowID, err := bot.FindAnswerRowID(event.Source.UserID)
				if err != nil {
					log.Fatal(err)
				}
				if answerRowID != 0 {
					message, err := bot.SaveAnswerAndGetNextMessage("NULL", answerRowID, event.Source.UserID)
					if err != nil {
						log.Fatal(err)
					}

					bot.ReplyMessage(event.ReplyToken, message).Do()
					return
				}
			}
		}
	}
}
