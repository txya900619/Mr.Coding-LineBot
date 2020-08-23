package spreadsheets

import (
	"errors"
	"fmt"
	"strconv"
)

func (ss *Spreadsheets) FindAnswerRowID(userID string) (int, error) {
	// userID col range in spreadsheets
	userIdRange := "I:I"

	resp, err := ss.Values.Get(ss.SpreadsheetsID, userIdRange).Do()
	if err != nil {
		return 0, fmt.Errorf("get spreadsheets value fail, err: %v", err)
	}

	if len(resp.Values) == 0 {
		return 0, errors.New("no data found")
	}

	for index, row := range resp.Values {
		if row[0] == userID {
			// index+1 will be rowID
			return index + 1, nil
		}
	}

	// if userID no match any row
	return 0, nil
}

func (ss *Spreadsheets) FindCurrentQuestionColID(answerRowID int) (ColumnID, error) {
	answerRowIDStr := strconv.Itoa(answerRowID)
	answerRange := answerRowIDStr + ":" + answerRowIDStr
	resp, err := ss.Values.Get(ss.SpreadsheetsID, answerRange).Do()
	if err != nil {
		return 0, fmt.Errorf("get spreadsheets value fail, err: %v", err)
	}

	for index, col := range resp.Values[0] {
		if col == "" {
			return ColumnID(index), nil
		}
	}

	// if all answer is fill
	return 0, nil
}

func (ss *Spreadsheets) GetLastRowID() (int, error) {
	timestampRange := "A:A"
	resp, err := ss.Values.Get(ss.SpreadsheetsID, timestampRange).Do()
	if err != nil {
		return 0, fmt.Errorf("get spreadsheets value fail, err: %v", err)
	}

	return len(resp.Values) + 1, nil

}
