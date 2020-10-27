package drive

import (
	"fmt"
	"io"

	"google.golang.org/api/drive/v3"
)

//UploadNewFile is function that can upload file(img) to google drive
func (d *Drive) UploadNewFile(data io.ReadCloser, name string) (string, error) {
	newFile := &drive.File{Name: name, Parents: []string{d.FolderID}}
	file, err := d.Files.Create(newFile).Media(data).Do()
	if err != nil {
		return "", fmt.Errorf("upload file fail, err: %v", err)
	}
	return "https://drive.google.com/open?id=" + file.Id, nil
}
