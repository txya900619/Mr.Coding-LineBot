package drive

import (
	"context"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"io/ioutil"
)

type Drive struct {
	*drive.Service
	FolderID string
}

func New(folderID string) (*Drive, error) {
	b, err := ioutil.ReadFile("token.json")

	if err != nil {

		return nil, fmt.Errorf("can't read token.json, err: %v", err)
	}

	token, err := google.JWTConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		return nil, fmt.Errorf("get config fail, err: %v", err)
	}

	client := token.Client(context.Background())

	service, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Drive client, err: %v", err)
	}

	drive := &Drive{Service: service, FolderID: folderID}

	return drive, err
}
