package slack

import (
	"fmt"
	"github.com/slack-go/slack"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type conf struct {
	Token   string `yaml:"slack_api_token"`
	Channel string `yaml:"slack_channel"`
}

func UploadFile() error {
	c, err := GetCredentials()

	if err != nil {
		return err
	}

	api := slack.New(c.Token)

	params := slack.FileUploadParameters{
		Title:          "OlTP ",
		Filetype:       "json",
		File:           "report/sample/sample_oltp.json",
		Channels:       []string{c.Channel},
		InitialComment: "This is sample OLTP",
	}

	_, err = api.UploadFile(params)
	if err != nil {
		return err
	}

	return nil

}

func GetCredentials() (*conf, error) {

	//Read from the config file

	var c conf
	fmt.Println(os.Getwd())
	yamlFile, err := ioutil.ReadFile("config/config.yaml")
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil

}
