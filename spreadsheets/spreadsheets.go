package spreadsheets

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"io/ioutil"
	"strconv"
)

type ColumnID rune

const (
	TimeStamp ColumnID = iota
	QuestionEmail
	QuestionName
	QuestionStudentNo
	QuestionDepartment
	QuestionProgramming
	QuestionUploadFile
	QuestionNote
	UserID
)

type Spreadsheets struct {
	*sheets.SpreadsheetsService
	SpreadsheetsID string
}

func (colID *ColumnID) String() string {
	return string(rune(*colID) + 'A')
}

func New(spreadsheetsID string) (*Spreadsheets, error) {
	b, err := ioutil.ReadFile("token.json")

	if err != nil {

		return nil, fmt.Errorf("can't read token.json, err: %v", err)
	}

	token, err := google.JWTConfigFromJSON(b, sheets.SpreadsheetsScope)
	if err != nil {
		return nil, fmt.Errorf("get config fail, err: %v", err)
	}

	client := token.Client(context.Background())

	service, err := sheets.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Drive client, err: %v", err)
	}

	spreadsheets := &Spreadsheets{SpreadsheetsService: service.Spreadsheets, SpreadsheetsID: spreadsheetsID}

	return spreadsheets, err
}

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

func (ss *Spreadsheets) SaveValueToSpecificCell(value, ranges string) error {
	_, err := ss.Values.Update(ss.SpreadsheetsID, ranges, &sheets.ValueRange{Values: [][]interface{}{{value}}}).ValueInputOption("RAW").Do()
	if err != nil {
		return err
	}
	return nil
}
