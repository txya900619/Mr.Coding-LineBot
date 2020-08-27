package entroy

import (
	"math/rand"
	"strconv"
)

func (p *Player) RandomQuestionID() uint {

	// How many Question in questions.json
	questionID := uint(rand.Intn(20))
	for !p.checkQuestionDuplicate(questionID) {
		questionID = uint(rand.Intn(20))
	}

	return questionID
}

// If duplicate return false
func (p *Player) checkQuestionDuplicate(questionID uint) bool {
	for _, answered := range p.AnsweredList {
		if answered == questionID {
			return false
		}
	}
	return true
}

func CheckAnswerText(answer string) bool {
	result := false
	for i := 1; i < 5; i++ {
		if strconv.Itoa(i) == answer {
			result = true
		}
	}
	return result
}
