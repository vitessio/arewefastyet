package slack

import (
	"errors"
	"github.com/slack-go/slack"
	"path"
)

type (
	FileUploadMessage struct {
		Title string
		Comment string
		FilePath string
		FileType string
	}

	TextMessage struct {
		Content string
	}

	Message interface {
		Send(config Config) error
	}
)

func (f FileUploadMessage) Send(config Config) (err error) {
	if !config.IsValid() {
		return errors.New(ErrorInvalidConfiguration)
	}

	api := slack.New(config.Token)

	if f.FileType == "" {
		getFileType(&f)
	}

	params := slack.FileUploadParameters{
		Title:          f.Title,
		Filetype:       f.FileType,
		File:           f.FilePath,
		Channels:       []string{config.Channel},
		InitialComment: f.Comment,
	}

	_, err = api.UploadFile(params)
	if err != nil {
		return err
	}
	return nil
}

func getFileType(f *FileUploadMessage) {
	ext := path.Ext(f.FilePath)
	if ext == "" {
		ext = ".txt"
	}
	f.FileType = ext[1:]
}

func (t TextMessage) Send(config Config) (err error) {
	api := slack.New(config.Token)

	_, _, err = api.PostMessage(config.Channel, slack.MsgOptionText(t.Content, false))

	if err != nil {
		return err
	}
	return nil
}
