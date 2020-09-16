package main

import (
	"Mr.Coding-LineBot/config"
	"Mr.Coding-LineBot/entroy"
	"Mr.Coding-LineBot/mrcoding"
	"Mr.Coding-LineBot/spreadsheets"
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

						currentQuestion := entroy.Questions[player.AnsweredList[len(player.AnsweredList)-1]]

						if len(player.AnsweredList) == 5 {
							currentQuestion.Final = true
						}

						if message.Text == currentQuestion.Answer {
							player.Score++
						}

						entroyBot.ReplyMessage(event.ReplyToken, currentQuestion.ReasonMessage(message.Text)).Do()
						return

					} else if message.Text == "社團博覽會有獎徵答" {
						entroyBot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(`不要再輸入 社團博覽會有獎徵答 拉 > <`)).Do()
						return
					} else {
						entroyBot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("要輸入 1 到 4 歐")).Do()
						return
					}
				default:
					entroyBot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("不要傳文字以外的訊息歐 ～～")).Do()
					return
				}
			case linebot.EventTypePostback:
				switch event.Postback.Data {
				case "next":
					player := entroyBot.PlayerList[event.Source.UserID]
					if len(player.AnsweredList) >= 5 {
						message := player.FinalMessage()
						delete(entroyBot.PlayerList, event.Source.UserID)

						entroyBot.ReplyMessage(event.ReplyToken, message).Do()
						return
					} else {
						nextQuestionID := player.RandomQuestionID()
						message := player.GetQuestionMessageByID(nextQuestionID)
						player.AnsweredList = append(player.AnsweredList, nextQuestionID)

						entroyBot.ReplyMessage(event.ReplyToken, message).Do()
						return
					}
				}
			}
		} else {
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

						entroyBot.ReplyMessage(event.ReplyToken, entroy.StartMessage()).Do()
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
							return
						} else {
							bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("點開選單選擇功能，\n或輸入 /help 選擇想使用的功能。")).Do()
							return
						}
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
						fileUrl, err := bot.Drive.UploadNewFile(content.Content, event.Timestamp.String()+"-"+event.Source.UserID)
						if err != nil {
							log.Fatal(err)
						}

						message, err := bot.SaveAnswerAndGetNextMessage(fileUrl, answerRowID, event.Source.UserID)
						if err != nil {
							log.Fatal(err)
						}

						bot.ReplyMessage(event.ReplyToken, message).Do()
						return
					} else {
						bot.ReplyMessage(event.ReplyToken, linebot.NewFlexMessage("error", mrcoding.GetTypeErrorFlexContainer())).Do()
						return
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
