package slack

import (
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	configPath     = "../../config/config.yaml"
	reportFilePath = "../../report/sample/sample_oltp.json"
)

func TestReadCred(t *testing.T) {

	c, err := GetCredentials(configPath)

	require.NoError(t, err)
	require.NotEmpty(t, c.Token)
	require.NotEmpty(t, c.Channel)
}

func TestUploadFile(t *testing.T) {
	require.NoError(t, UploadFile(configPath, reportFilePath))
}
