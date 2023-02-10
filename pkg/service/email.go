package service

import (
	"fmt"
	"os"

	"github.com/String-xyz/go-lib/common"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendEmail(from string, to string, fromAddress string, toAddress string, subject string, body string) error {
	if fromAddress != "auth@string.xyz" {
		fromAddress = "auth@string.xyz" // TODO: create a new sender
	}
	f := mail.NewEmail(from, fromAddress)
	t := mail.NewEmail(to, toAddress)
	textContent := ""
	// TODO: Parse body to ensure links are valid
	message := mail.NewSingleEmail(f, subject, t, textContent, body)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	res, err := client.Send(message)

	if res.StatusCode >= 400 {
		fmt.Println(">>> email failed to send: ", res.Body)
		return common.StringError(errors.New(res.Body))
	}

	if err != nil {
		return common.StringError(err)
	}
	return nil
}
