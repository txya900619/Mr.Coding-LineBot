package main

import (
	"Mr.Coding-LineBot/config"
	"Mr.Coding-LineBot/entroy"
	"Mr.Coding-LineBot/mrcoding"
	"fmt"
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
	"net/http"
)

var bot *mrcoding.Bot
var entroyBot *entroy.Bot

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
	entroyBot = entroy.New(bot.Client)

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

		if entroyBot.CheckPlayerInGame(event.Source.UserID) {
			switch event.Type {
			case linebot.EventTypeMessage:
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					if entroy.CheckAnswerText(message.Text) {
						player := entroyBot.PlayerList[event.Source.UserID]

						messages := make([]linebot.SendingMessage, 0)
						currentQuestion := entroy.Questions[player.AnsweredList[len(player.AnsweredList)-1]]
						messages = append(messages, currentQuestion.ReasonMessage(message.Text))

						if message.Text == currentQuestion.Answer {
							player.Score++
						}

						if len(player.AnsweredList) >= 5 {
							messages = append(messages, player.FinalMessage())
							delete(entroyBot.PlayerList, event.Source.UserID)

							entroyBot.ReplyMessage(event.ReplyToken, messages...).Do()
							return
						} else {
							nextQuestionID := player.RandomQuestionID()
							player.AnsweredList = append(player.AnsweredList, nextQuestionID)
							messages = append(messages, player.GetQuestionMessageByID(nextQuestionID))

							entroyBot.ReplyMessage(event.ReplyToken, messages...).Do()
							return
						}
					} else {
						entroyBot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("要輸入 1 到 4 歐")).Do()
						return
					}
				default:
					entroyBot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("不要傳文字以外的訊息歐 ～～")).Do()
					return
				}
			}
		} else {
			switch event.Type {
			case linebot.EventTypeMessage:
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					switch message.Text {
					case "Mr.Coding 表單":
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
							return
						} else {
							message, err := bot.SaveAnswerAndGetNextMessage(message.Text, answerRowID, event.Source.UserID)
							if err != nil {
								log.Fatal(err)
							}

							bot.ReplyMessage(event.ReplyToken, message).Do()
							return
						}
					case "社團博覽會有獎徵答":
						player := &entroy.Player{
							AnsweredList: make([]uint, 0),
							Score:        0,
						}
						entroyBot.PlayerList[event.Source.UserID] = player

						// How many Question in questions.json
						questionID := player.RandomQuestionID()
						player.AnsweredList = append(player.AnsweredList, questionID)

						entroyBot.ReplyMessage(event.ReplyToken, entroy.StartMessage(), player.GetQuestionMessageByID(questionID)).Do()
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
							return
						}
					}
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
}
