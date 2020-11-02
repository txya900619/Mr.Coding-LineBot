package mrcoding

import (
	"errors"
	"fmt"
	"log"
	"time"

	"Mr.Coding-LineBot/config"
	"Mr.Coding-LineBot/drive"
	"Mr.Coding-LineBot/mrcoding/messages"
	"Mr.Coding-LineBot/spreadsheets"
	"github.com/go-playground/validator/v10"
	"github.com/gomodule/redigo/redis"
	"github.com/line/line-bot-sdk-go/linebot"
)

//Bot is a struct that and handle all need of linebot using
type Bot struct {
	//Origin line-sdk linebot client
	*linebot.Client

	//An google spreadsheets client for write and read spreadsheets
	Spreadsheets *spreadsheets.Spreadsheets

	//An google diver client for upload picture to diver
	Drive *drive.Drive

	//An token that can authorization when send request to Mr.Coding backend
	backendToken string

	//An redis conn, using write and read redis's data
	Redis redis.Conn
}

//ErrEmailValidateFail is a const error fot email validator
var ErrEmailValidateFail = errors.New("email validator fail")

//New can create an new Bot instance
func New(c *config.Config, options ...linebot.ClientOption) (*Bot, error) {
	//create spreadsheets client instance
	ss, err := spreadsheets.New(c.SpreadsheetId)
	if err != nil {
		return nil, err
	}

	//create drive client instance
	drive, err := drive.New(c.FolderId)

	//create origin linebot client instance
	lb, err := linebot.New(c.ChannelSecret, c.ChannelToken, options...)
	if err != nil {
		return nil, err
	}

	//create redis client instace
	redis, err := redis.DialURL("redis://redis:6379")
	if err != nil {
		return nil, err
	}

	return &Bot{lb, ss, drive, c.CreateChatroomToken, redis}, nil
}

//Initail questions by save some user data
func (bot *Bot) initQuestions(userID string) (*linebot.FlexMessage, error) {
	// save answer to index 0
	bot.Redis.Do("RPUSH", userID+"Data", time.Now().String())

	// save what question should be answered
	bot.Redis.Do("SET", userID, string(rune(spreadsheets.QuestionEmail)))

	//get first question message (email)
	message := messages.QuestionMessage(spreadsheets.QuestionEmail)

	return message, nil
}

//endQuestions is an function to do some thing that need to do after last answer save to redis
//like send answers to spreadsheets and backend or delete answers in redis
func (bot *Bot) endQuestions(userID string) string {
	//get answers from redis list
	row, err := redis.Strings(bot.Redis.Do("LRANGE", userID+"Data", 0, 7))
	if err != nil {
		log.Fatalf("redis LRANGE answers fail, err: %v", err)
	}

	//append answers to spreadsheets
	err = bot.Spreadsheets.AppendRow(row)
	if err != nil {
		log.Fatal(err)
	}

	//delete answers in redis list
	deleteUserAnswer(userID, bot.Redis)

	//send request to backend and return message that have chatroom's URL
	return bot.createChatroomAndGetID(userID, row[2]+"的詢問聊天室")
}

//Save answer and get next question
func (bot *Bot) saveAnswerAndGetNextQuestion(currentPosition, answer, userID string) (*linebot.FlexMessage, error) {
	//transform currentPostion(string) to ColumnID
	questionColID := spreadsheets.ColumnID([]rune(currentPosition)[0])

	//save the answer to current question
	err := saveAnswer(answer, userID, questionColID, bot.Redis)
	if err != nil {
		//if email validate fail then return email error message to re-question
		if err == ErrEmailValidateFail {
			return messages.EmailErrorMessage(), nil

		}
		return nil, err
	}

	//if current question is last question
	if questionColID == spreadsheets.QuestionNote {
		//do something need to do in last
		chatroomID := bot.endQuestions(userID)

		return messages.CompleteFormMessage(chatroomID), nil
	}

	//update currentPostion
	_, err = bot.Redis.Do("SET", userID, string(rune(questionColID)+1))
	if err != nil {
		return nil, fmt.Errorf("redis SET current position to next fail, err: %v", err)
	}

	// return next question
	return messages.QuestionMessage(spreadsheets.ColumnID(rune(questionColID) + 1)), nil

}

//Save answer to redis
func saveAnswer(answer, userID string, questionColID spreadsheets.ColumnID, conn redis.Conn) error {
	//if question to answer is email question, then validate answer, if validate fail then re-question
	if questionColID == spreadsheets.QuestionEmail {
		v := validator.New()
		err := v.Var(answer, "email")
		if err != nil {
			return ErrEmailValidateFail
		}

	}

	//save answer to redis list
	_, err := conn.Do("RPUSH", userID+"Data", answer)
	if err != nil {
		return fmt.Errorf("redis RPUSH answer fail, err: %v", err)
	}

	return nil
}

//Delete all user data (answers and info) after answers is been send to spreadsheets and backend
func deleteUserAnswer(userID string, conn redis.Conn) error {
	//delete current position
	_, err := conn.Do("DEL", userID)
	if err != nil {
		return fmt.Errorf("redis DEL current position fail, err: %v", err)
	}

	//delete answers
	_, err = conn.Do("DEL", userID+"Data")
	if err != nil {
		return fmt.Errorf("redis DEL answers fail, err: %v", err)
	}

	return nil
}
