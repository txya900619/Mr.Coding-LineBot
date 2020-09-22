package main

import (
	"fmt"
	"log"
	"net/http"

	"Mr.Coding-LineBot/config"
	"Mr.Coding-LineBot/mrcoding"
	"Mr.Coding-LineBot/spreadsheets"
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
		switch event.Type {
		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				switch message.Text {
				case "Mr.Coding 表單":
					//bot.ReplyMessage(event.ReplyToken, linebot.NewFlexMessage("working", mrcoding.GetWorkingFlexContainer())).Do()
					//return
					answerRowID, err := bot.Spreadsheets.FindAnswerRowID(event.Source.UserID)
					if err != nil {
						log.Fatal(err)
					}
					if answerRowID == 0 {
						message, err := bot.QuestionStart(event.Source.UserID)
						if err != nil {
							log.Fatal(err)
						}

						bot.ReplyMessage(event.ReplyToken, message).Do()
					} else {
						message, err := bot.SaveAnswerAndGetNextMessage(message.Text, answerRowID, event.Source.UserID)
						if err != nil {
							log.Fatal(err)
						}

						bot.ReplyMessage(event.ReplyToken, message).Do()
					}
					return
				case "社團博覽會有獎徵答":

					bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("社團博覽會已結束")).Do()
					return
				case "/help":
					bot.ReplyMessage(event.ReplyToken, mrcoding.HelpMessage()).Do()
					return
				default:
					answerRowID, err := bot.Spreadsheets.FindAnswerRowID(event.Source.UserID)
					if err != nil {
						log.Fatal(err)
					}

					if answerRowID != 0 {
						message, err := bot.SaveAnswerAndGetNextMessage(message.Text, answerRowID, event.Source.UserID)
						if err != nil {
							log.Fatal(err)
						}

						bot.ReplyMessage(event.ReplyToken, message).Do()
					} else {
						bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("點開選單選擇功能，\n或輸入 /help 選擇想使用的功能。")).Do()
					}
					return
				}
			case *linebot.ImageMessage:
				answerRowID, err := bot.Spreadsheets.FindAnswerRowID(event.Source.UserID)
				if err != nil {
					log.Fatal(err)
				}
				questionColID, err := bot.Spreadsheets.FindCurrentQuestionColID(answerRowID)
				if err != nil {
					log.Fatal(err)
				}
				if questionColID == spreadsheets.QuestionUploadFile {
					content, err := bot.GetMessageContent(message.ID).Do()
					if err != nil {
						log.Fatal(err)
					}
					fileURL, err := bot.Drive.UploadNewFile(content.Content, event.Timestamp.String()+"-"+event.Source.UserID)
					if err != nil {
						log.Fatal(err)
					}

					message, err := bot.SaveAnswerAndGetNextMessage(fileURL, answerRowID, event.Source.UserID)
					if err != nil {
						log.Fatal(err)
					}

					bot.ReplyMessage(event.ReplyToken, message).Do()
				} else {
					bot.ReplyMessage(event.ReplyToken, linebot.NewFlexMessage("error", mrcoding.GetTypeErrorFlexContainer())).Do()
				}
				return

			}
		case linebot.EventTypePostback:
			switch event.Postback.Data {
			case "pass":
				answerRowID, err := bot.Spreadsheets.FindAnswerRowID(event.Source.UserID)
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
