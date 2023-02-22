package mail

import (
	"github.com/MathPeixoto/go-financial-system/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	config, err := util.LoadConfig("../app.env")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	subject := "A test email"
	content := `
		<h1>Hi there!</h1>
		<p>This is a test email</p>
	`
	to := []string{"matheus-pxt@hotmail.com"}
	attachFiles := []string{"../README.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
