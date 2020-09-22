package spreadsheets

import (
	"context"
	"fmt"
	"io/ioutil"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type ColumnID rune

const (
	TimeStamp ColumnID = iota + 'A'
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

	fmt.Println(b)

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
