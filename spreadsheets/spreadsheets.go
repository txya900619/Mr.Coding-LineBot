package spreadsheets

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"io/ioutil"
	"log"
	"os"
)

func New() (*sheets.Service, error) {
	b, err := ioutil.ReadFile("credentials.json")

	if err != nil {

		return nil, fmt.Errorf("can't read credentials.json, err: %v", err)
	}

	config, err := google.ConfigFromJSON(b, sheets.SpreadsheetsScope)

	if err != nil {
		return nil, fmt.Errorf("get config fail, err: %v", err)
	}

	token := getToken(config)

	service, err := sheets.NewService(context.Background(), option.WithTokenSource(config.TokenSource(context.Background(), token)))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Drive client, err: %v", err)
	}

	return service, err
}

func getToken(config *oauth2.Config) *oauth2.Token {
	token, err := tokenFromFile("token.json")
	if err != nil {
		token = getTokenFromWeb(config)
		saveToken("token.json", token)
	}

	return token
}

func tokenFromFile(filePath string) (*oauth2.Token, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	token := &oauth2.Token{}

	err = json.NewDecoder(f).Decode(token)

	return token, nil
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Can't read code %v", authCode)
	}

	token, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}

	return token
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
