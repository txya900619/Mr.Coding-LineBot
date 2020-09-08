package drive

import (
	"google.golang.org/api/drive/v3"
	"io"
)

func (d *Drive) UploadNewFile(data io.ReadCloser, name string) (string, error) {
	newFile := &drive.File{Name: name, Parents: []string{d.FolderID}}
	file, err := d.Files.Create(newFile).Media(data).Do()
	return "https://drive.google.com/open?id=" + file.Id, err
}
