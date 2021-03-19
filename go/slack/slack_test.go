package slack

import (
	"fmt"
	qt "github.com/frankban/quicktest"
	"os"
	"testing"
)

func TestReadCred(t *testing.T) {

	tests := []struct {
		name    string
		content map[string]string
		wantFile bool
		wantErr bool
	}{
		{name: "Test valid slack credentials", content: map[string]string{"slack_api_token": "xoxb-2384bizzoEHwq4GOykIR", "slack_channel": "test-channel"}, wantFile: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)
			path := c.TempDir()
			configPath := path + "/" + "config.yaml"
			file, err := os.Create(path + "/" + "config.yaml")
			if err != nil {
				c.Skip(err.Error())
			}

			for key, value := range tt.content {
				fmt.Fprintf(file, "%s: %s\n", key, value)
			}

			cred, err := GetCredentials(configPath)
			c.Assert(err, qt.IsNil)
			c.Assert(cred.Token, qt.HasLen, len(tt.content["slack_api_token"]))
			c.Assert(cred.Channel, qt.HasLen, len(tt.content["slack_channel"]))
		})
	}
}

func TestUploadFile(t *testing.T) {

	tests := []struct {
		name    string
		credentials map[string]string
		test_report string
		wantFile bool
		wantErr bool
	}{
		{name: "Test upload file", credentials: map[string]string{"slack_api_token": "wrong token", "slack_channel": "wrong channel"},
			test_report: "slack bot did you work", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)
			path := c.TempDir()
			configPath := path + "/" + "config.yaml"
			file, err := os.Create(path + "/" + "config.yaml")
			if err != nil {
				c.Skip(err.Error())
			}

			for key, value := range tt.credentials {
				fmt.Fprintf(file, "%s: %s\n", key, value)
			}
			reportPath := path + "/" + "report.json"
			fileReport, err := os.Create(path + "/" + "report.json")
			if err != nil {
				c.Skip(err.Error())
			}
			fmt.Fprintf(fileReport, "%s\n", tt.test_report)

			err = UploadFile(configPath, reportPath)
			if tt.wantErr == true {
				// Looks like you gave the right key (You special human)
				c.Assert(err, qt.Not(qt.IsNil))
			} else {
				c.Assert(err, qt.IsNil)
			}
		})
	}
}
