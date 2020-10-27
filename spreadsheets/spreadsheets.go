package spreadsheets

import (
	"context"
	"fmt"
	"io/ioutil"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

//ColumnID is alias of rune, using to tag what question it is
type ColumnID rune

const (
	//TimeStamp is time of question answered time, it will auto generate, don't need answer
	TimeStamp ColumnID = iota + '0'

	//QuestionEmail is question that user should answer it's email
	QuestionEmail

	//QuestionName is question that user should answer it's name
	QuestionName

	//QuestionStudentNo is question that user should answer it's student no.
	QuestionStudentNo

	//QuestionDepartment is question that user should answer it's department(系級)
	QuestionDepartment

	//QuestionProgramming is question that user should answer it's programming problem
	QuestionProgramming

	//QuestionUploadFile is question that user should answer it's programming problem shotcut, but it not required
	QuestionUploadFile

	//QuestionNote is question that user should answer it's note(any), but it not required
	QuestionNote
)

//Spreadsheets is a struct that contain origin spreadsheets client and SpreadsheetsID
type Spreadsheets struct {
	*sheets.SpreadsheetsService
	SpreadsheetsID string
}

func (colID *ColumnID) String() string {
	return string(rune(*colID) + 'A')
}

//New will read token from token.json and create new Spreadsheets instance
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
		return nil, fmt.Errorf("unable to retrieve Sheets client, err: %v", err)
	}

	spreadsheets := &Spreadsheets{SpreadsheetsService: service.Spreadsheets, SpreadsheetsID: spreadsheetsID}

	return spreadsheets, err
}
