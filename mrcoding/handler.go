package mrcoding

import (
	"log"
	"net/http"

	"Mr.Coding-LineBot/mrcoding/messages"
	"Mr.Coding-LineBot/spreadsheets"
	"github.com/gomodule/redigo/redis"
	"github.com/line/line-bot-sdk-go/linebot"
)

//Handler is function return an http handler function to serve linebot
func Handler(bot *Bot) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
						messageToSend = linebot.NewTextMessage("預計 11 月開放使用，目前竭力開發中～")
					case "/help":
						messageToSend = messages.HelpMessage()

					case "星爆氣流斬":
						currentPosition, err := redis.String(bot.Redis.Do("GET", event.Source.UserID))
						if err != nil {
							if err == redis.ErrNil {
								messageToSend, err = bot.initQuestions(event.Source.UserID)
								if err != nil {
									log.Fatal(err)
								}
								break
							} else {
								log.Fatalf("get current position fail, err: %v", err)
							}
						}

						messageToSend, err = bot.saveAnswerAndGetNextQuestion(currentPosition, message.Text, event.Source.UserID)
						if err != nil {
							log.Fatal(err)
						}
					default:
						currentPosition, err := redis.String(bot.Redis.Do("GET", event.Source.UserID))
						if err != nil {
							if err == redis.ErrNil {
								messageToSend = linebot.NewTextMessage("點開選單選擇功能，\n或輸入 /help 選擇想使用的功能。")
								break
							} else {
								log.Fatalf("get current position fail, err: %v", err)
							}
						}

						messageToSend, err = bot.saveAnswerAndGetNextQuestion(currentPosition, message.Text, event.Source.UserID)
						if err != nil {
							log.Fatal(err)
						}

					}
				case *linebot.ImageMessage:
					currentPosition, err := redis.String(bot.Redis.Do("GET", event.Source.UserID))

					if err != nil {
						if err == redis.ErrNil {
							messageToSend = messages.TypeErrorMessage()
						} else {
							log.Fatalf("get current position fail, err: %v", err)
						}
					}

					if spreadsheets.ColumnID([]rune(currentPosition)[0]) == spreadsheets.QuestionUploadFile {
						content, err := bot.GetMessageContent(message.ID).Do()
						if err != nil {
							log.Fatal(err)
						}
						fileURL, err := bot.Drive.UploadNewFile(content.Content, event.Timestamp.String()+"-"+event.Source.UserID)
						if err != nil {
							log.Fatal(err)
						}

						messageToSend, err = bot.saveAnswerAndGetNextQuestion(currentPosition, fileURL, event.Source.UserID)
						if err != nil {
							log.Fatal(err)
						}
					} else {
						messageToSend = messages.TypeErrorMessage()
					}

				}
			case linebot.EventTypePostback:
				switch event.Postback.Data {
				case "pass":
					currentPosition, err := redis.String(bot.Redis.Do("GET", event.Source.UserID))
					if err != nil {
						if err != redis.ErrNil {
							log.Fatalf("get current position fail, err: %v", err)
						}
					}

					messageToSend, err = bot.saveAnswerAndGetNextQuestion(currentPosition, "NULL", event.Source.UserID)
					if err != nil {
						log.Fatal(err)
					}
				}
			}

			bot.ReplyMessage(event.ReplyToken, messageToSend).Do()
		}
	}

}
