package main

import (
	"Mr.Coding-LineBot/config"
	"Mr.Coding-LineBot/mrcoding"
	"Mr.Coding-LineBot/spreadsheets"
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
	"net/http"
)

var bot *mrcoding.Bot
var ss *spreadsheets.Spreadsheets

func main() {
	// Get config at ./config.yaml
	c, err := config.New()
	if err != nil {
		log.Fatalf("Read config.yaml file fail, %v", err)
	}

	ss, err = spreadsheets.New(c.SpreadsheetId)
	if err != nil {
		log.Fatal(err)
	}

	bot, err = mrcoding.New(c.ChannelSecret, c.ChannelToken, ss)
	if err != nil {
		log.Fatalf("Create linebot fail, %v", err)
	}

	err = bot.SaveValueToSpecificCell("cc", "A18:A18")
	if err != nil {
		log.Fatal(err)
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
				switch message.Text {
				case "Mr.Coding 表單":
					return
				case "社團博覽會有獎徵答":
					return
				default:
					answerRowID, err := ss.FindAnswerRowID(event.Source.UserID)
					if err != nil {
						log.Fatal(err)
					}

					if answerRowID != 0 {
						question, err := bot.SaveAnswerAndGetNextQuestion(message.Text, answerRowID)
						if err != nil {
							log.Fatal(err)
						}

						bot.ReplyMessage(event.ReplyToken, question)
						return
					}
				}

			}
		}
	}
}
