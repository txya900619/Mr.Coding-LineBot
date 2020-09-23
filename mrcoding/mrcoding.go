package mrcoding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"Mr.Coding-LineBot/config"
	"Mr.Coding-LineBot/drive"
	"Mr.Coding-LineBot/spreadsheets"
	"github.com/go-playground/validator/v10"
	"github.com/gomodule/redigo/redis"
	"github.com/line/line-bot-sdk-go/linebot"
)

type Bot struct {
	*linebot.Client
	Spreadsheets *spreadsheets.Spreadsheets
	Drive        *drive.Drive
	backendToken string
	Redis        redis.Conn
}

func New(c *config.Config, options ...linebot.ClientOption) (*Bot, error) {
	ss, err := spreadsheets.New(c.SpreadsheetId)
	if err != nil {
		return nil, err
	}

	drive, err := drive.New(c.FolderId)

	lb, err := linebot.New(c.ChannelSecret, c.ChannelToken, options...)
	if err != nil {
		return nil, err
	}

	redis, err := redis.DialURL("redis://redis:6379")
	if err != nil {
		return nil, err
	}

	return &Bot{lb, ss, drive, c.CreateChatroomToken, redis}, nil
}

func (bot *Bot) QuestionStart(userID string) (*linebot.FlexMessage, error) {
	// save answer to index 0
	bot.Redis.Do("ZADD", userID+"Data", 0, time.Now().String())

	// save what question should be answered
	bot.Redis.Do("SET", userID, string(rune(spreadsheets.QuestionEmail)))

	flexContainer := getQuestionFlexContainer(spreadsheets.QuestionEmail)
	message := linebot.NewFlexMessage("Questions", flexContainer)
	return message, nil
}

func (bot *Bot) SaveAnswerAndGetNextMessage(answer string, currentPosition string, userID string) (*linebot.FlexMessage, error) {
	questionColID := spreadsheets.ColumnID([]rune(currentPosition)[0])

	if questionColID == spreadsheets.QuestionEmail {
		v := validator.New()
		err := v.Var(answer, "email")
		if err != nil {
			return linebot.NewFlexMessage("email input error", getEmailErrorFlexContainer()), nil
		}
	}

	bot.Redis.Do("ZADD", userID+"Data", string(rune(questionColID)), answer)
	// If is last question
	if questionColID == spreadsheets.QuestionNote {
		_, err := bot.Redis.Do("DEL", userID)
		if err != nil {
			return nil, err
		}

		row, err := redis.MultiBulk(bot.Redis.Do("ZRANGE", userID+"Data", 0, 7))
		rowStr, err := redis.Strings(bot.Redis.Do("ZRANGE", userID+"Data", 0, 7))
		if err != nil {
			return nil, err
		}
		fmt.Println(rowStr)

		err = bot.Spreadsheets.AppendRow(row)
		if err != nil {
			return nil, err
		}

		bot.Redis.Do("DEL", userID+"Data")
		chatroomID := bot.createChatroomAndGetID(userID)
		flexContainer := getCompleteFormFlexContainer(chatroomID)
		message := linebot.NewFlexMessage("Final", flexContainer)
		return message, nil
	}

	bot.Redis.Do("SET", userID, string(rune(questionColID)+1))
	flexContainer := getQuestionFlexContainer(spreadsheets.ColumnID(rune(questionColID) + 1))
	message := linebot.NewFlexMessage("Questions", flexContainer)
	return message, nil
}

func (bot *Bot) createChatroomAndGetID(userID string) string {
	client := &http.Client{}
	fmt.Println(userID)
	reqBody := map[string]string{"lineChatroomUserID": userID}
	jsonReqBody, _ := json.Marshal(reqBody)
	req, err := http.NewRequest(http.MethodPost, "https://mrcoding.org/api/chatrooms", bytes.NewBuffer(jsonReqBody))
	if err != nil {
		log.Fatalf("newReq, err: %v", err)
	}

	req.Header.Set("Authorization", bot.backendToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("reqDo, err: %v", err)
	}

	//Backend err
	if !(resp.StatusCode == 200 || resp.StatusCode == 201) {
		log.Fatal("statusCode err")

	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Fatalf("readAll, err: %v", err)
	}

	result := make(map[string]interface{})
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		log.Fatalf("jsonUnmarshal, err: %v", err)
	}

	return result["_id"].(string)
}
