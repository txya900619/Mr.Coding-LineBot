package mrcoding

import (
	"Mr.Coding-LineBot/config"
	"Mr.Coding-LineBot/spreadsheets"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/line/line-bot-sdk-go/linebot"
	"io/ioutil"
	"log"
	"net/http"
)

type Bot struct {
	*linebot.Client
	Spreadsheets *spreadsheets.Spreadsheets
	backendToken string
}

func New(c *config.Config, options ...linebot.ClientOption) (*Bot, error) {
	ss, err := spreadsheets.New(c.SpreadsheetId)
	if err != nil {
		return nil, err
	}

	lb, err := linebot.New(c.ChannelSecret, c.ChannelToken, options...)
	if err != nil {
		return nil, err
	}

	return &Bot{lb, ss, c.CreateChatroomToken}, nil
}

func (bot *Bot) QuestionStart(userID string) (*linebot.FlexMessage, error) {
	lastRowID, err := bot.Spreadsheets.GetLastRowID()
	if err != nil {
		return nil, err
	}

	err = bot.Spreadsheets.InsertTimestampAndUserID(userID, lastRowID)
	if err != nil {
		return nil, err
	}

	flexContainer := getQuestionFlexContainer(spreadsheets.QuestionEmail)
	message := linebot.NewFlexMessage("Questions", flexContainer)
	return message, nil
}

func (bot *Bot) SaveAnswerAndGetNextMessage(answer string, rowID int, userID string) (*linebot.FlexMessage, error) {
	questionColID, err := bot.Spreadsheets.FindCurrentQuestionColID(rowID)
	if err != nil {
		return nil, err
	}

	if questionColID == spreadsheets.QuestionEmail {
		v := validator.New()
		err := v.Var(answer, "email")
		if err != nil {
			return linebot.NewFlexMessage("email input error", getEmailErrorFlexContainer()), nil
		}
	}

	ranges := getRange(rowID, questionColID)

	err = bot.Spreadsheets.SaveValueToSpecificCell(answer, ranges)
	if err != nil {
		return nil, err
	}

	// If is last question
	if questionColID == spreadsheets.QuestionNote {
		err = bot.Spreadsheets.DeleteUserID(rowID)
		if err != nil {
			return nil, err
		}
		chatroomID := bot.createChatroomAndGetID(userID)
		flexContainer := getCompleteFormFlexContainer(chatroomID)
		message := linebot.NewFlexMessage("Final", flexContainer)
		return message, nil
	}

	flexContainer := getQuestionFlexContainer(spreadsheets.ColumnID(rune(questionColID) + 1))
	message := linebot.NewFlexMessage("Questions", flexContainer)
	return message, nil
}

func (bot *Bot) createChatroomAndGetID(userID string) string {
	client := &http.Client{}
	reqBodyBytes := []byte(fmt.Sprintf(`{"owner":"%v"}`, userID))
	req, err := http.NewRequest(http.MethodPost, "https://mrcoding.org/api/chatrooms", bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		log.Fatalf("newReq, err: %v", err)
	}

	req.Header.Set("Authorization", bot.backendToken)

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
