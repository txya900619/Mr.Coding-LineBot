package spreadsheets

import (
	"context"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"io/ioutil"
)

func New() (*sheets.Service, error) {
	b, err := ioutil.ReadFile("token.json")

	if err != nil {

		return nil, fmt.Errorf("can't read token.json, err: %v", err)
	}

	token, err := google.JWTConfigFromJSON(b, sheets.SpreadsheetsReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("get config fail, err: %v", err)
	}

	client := token.Client(context.Background())

	service, err := sheets.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Drive client, err: %v", err)
	}

	return service, err
}
