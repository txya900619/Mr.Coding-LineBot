package main

import (
	"fmt"
	"log"
	"net/http"

	"Mr.Coding-LineBot/config"
	"Mr.Coding-LineBot/mrcoding"
	"Mr.Coding-LineBot/spreadsheets"
	"github.com/gomodule/redigo/redis"
	"github.com/line/line-bot-sdk-go/linebot"
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
		var messageToSend linebot.SendingMessage
		switch event.Type {
		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				switch message.Text {
				case "Mr.Coding 表單":
					currentPosition, err := redis.String(bot.Redis.Do("GET", event.Source.UserID))

					if err == nil {
						messageToSend, err = bot.SaveAnswerAndGetNextMessage(message.Text, currentPosition, event.Source.UserID)
						if err != nil {
							log.Fatal(err)
						}
					} else if err == redis.ErrNil {
						messageToSend, err = bot.QuestionStart(event.Source.UserID)
						if err != nil {
							log.Fatal(err)
						}
					} else {
						log.Fatal(err)
					}

				case "社團博覽會有獎徵答":
					messageToSend = linebot.NewTextMessage("社團博覽會已結束")
				case "/help":
					messageToSend = mrcoding.HelpMessage()
				default:
					currentPosition, err := redis.String(bot.Redis.Do("GET", event.Source.UserID))
					if err == nil {
						messageToSend, err = bot.SaveAnswerAndGetNextMessage(message.Text, currentPosition, event.Source.UserID)
						if err != nil {
							log.Fatal(err)
						}
					} else if err == redis.ErrNil {
						messageToSend = linebot.NewTextMessage("點開選單選擇功能，\n或輸入 /help 選擇想使用的功能。")
					} else {
						log.Fatal(err)
					}
				}
			case *linebot.ImageMessage:
				currentPosition, err := redis.String(bot.Redis.Do("GET", event.Source.UserID))

				if err == nil {
					if spreadsheets.ColumnID([]rune(currentPosition)[0]) == spreadsheets.QuestionUploadFile {
						content, err := bot.GetMessageContent(message.ID).Do()
						if err != nil {
							log.Fatal(err)
						}
						fileURL, err := bot.Drive.UploadNewFile(content.Content, event.Timestamp.String()+"-"+event.Source.UserID)
						if err != nil {
							log.Fatal(err)
						}

						messageToSend, err = bot.SaveAnswerAndGetNextMessage(fileURL, currentPosition, event.Source.UserID)
						if err != nil {
							log.Fatal(err)
						}
					} else {
						messageToSend = linebot.NewFlexMessage("error", mrcoding.GetTypeErrorFlexContainer())
					}
				} else if err == redis.ErrNil {
					messageToSend = linebot.NewFlexMessage("error", mrcoding.GetTypeErrorFlexContainer())
				} else {
					log.Fatal(err)
				}

			}
		case linebot.EventTypePostback:
			switch event.Postback.Data {
			case "pass":
				currentPosition, err := redis.String(bot.Redis.Do("GET", event.Source.UserID))
				if err == nil {
					messageToSend, err = bot.SaveAnswerAndGetNextMessage("NULL", currentPosition, event.Source.UserID)
					if err != nil {
						log.Fatal(err)
					}
				} else if err != redis.ErrNil {
					log.Fatal(err)
				}
			}
		}

		bot.ReplyMessage(event.ReplyToken, messageToSend).Do()
	}
}
