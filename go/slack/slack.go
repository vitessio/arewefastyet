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

func UploadFile(configFilePath,reportFilePath string) error {
	c, err := GetCredentials(configFilePath)

	if err != nil {
		return err
	}

	api := slack.New(c.Token)

	params := slack.FileUploadParameters{
		Title:          "OlTP ",
		Filetype:       "json",
		File:           reportFilePath,
		Channels:       []string{c.Channel},
		InitialComment: "This is sample OLTP",
	}

	_, err = api.UploadFile(params)
	if err != nil {
		return err
	}

	return nil

}

func GetCredentials(configFilePath string) (*conf, error) {

	//Read from the config file

	var c conf
	fmt.Println(os.Getwd())
	yamlFile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil

}
