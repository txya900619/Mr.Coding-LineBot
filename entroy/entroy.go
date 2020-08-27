package entroy

import (
	"encoding/json"
	"github.com/line/line-bot-sdk-go/linebot"
	"io/ioutil"
	"log"
)

var Questions []Question

type Bot struct {
	*linebot.Client
	PlayerList map[string]*Player
}

type Player struct {
	AnsweredList []uint
	Score        uint
}

type Question struct {
	Question string   `json:"question"`
	Options  []string `json:"options"`
	Answer   string   `json:"answer"`
	Reason   string   `json:"reason"`
}

func New(client *linebot.Client) *Bot {
	rawByte, err := ioutil.ReadFile("questions.json")
	if err != nil {
		log.Fatalf("read Questions.json fail, err: %v", err)
	}

	result := make([]Question, 0)

	err = json.Unmarshal(rawByte, &result)
	if err != nil {
		log.Fatalf("Questions.json unmarshal fail, err: %v", err)
	}

	Questions = result

	return &Bot{client, make(map[string]*Player)}

}

// If player userID in player list return true
func (b *Bot) CheckPlayerInGame(userID string) bool {
	if b.PlayerList[userID] != nil {
		return true
	}
	return false
}
