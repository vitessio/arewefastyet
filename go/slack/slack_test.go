package slack

import (
	qt "github.com/frankban/quicktest"
	"testing"
)

func Test_getFileType(t *testing.T) {
	tests := []struct {
		name string
		f *FileUploadMessage
		want string
	}{
		{name: "Regular extension", f: &FileUploadMessage{FilePath: "file.json"}, want: "json"},
		{name: "No extension", f: &FileUploadMessage{FilePath: "file"}, want: "txt"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)

			getFileType(tt.f)
			c.Assert(tt.f.FileType, qt.Equals, tt.want)
		})
	}
}

func TestFileUploadMessage_Send(t *testing.T) {
	type fields struct {
		Title    string
		Comment  string
		FilePath string
		FileType string
	}
	type args struct {
		config Config
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "Invalid configuration (1)", args: args{config: Config{}}, wantErr: true},
		{name: "Invalid configuration (2)", args: args{config: Config{Token: "token"}}, wantErr: true},
		{name: "Invalid configuration (3)", args: args{config: Config{Channel: "channel"}}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FileUploadMessage{
				Title:    tt.fields.Title,
				Comment:  tt.fields.Comment,
				FilePath: tt.fields.FilePath,
				FileType: tt.fields.FileType,
			}
			c := qt.New(t)

			err := f.Send(tt.args.config)
			if tt.wantErr {
				c.Assert(err, qt.Not(qt.IsNil))
				return
			}
		})
	}
}