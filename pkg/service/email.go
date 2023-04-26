package service

import (
	"errors"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/platform-admin-api/config"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendEmail(from string, to string, toAddress string, subject string, body string) error {
	fromAddress := config.Var.AUTH_EMAIL_ADDRESS

	f := mail.NewEmail(from, fromAddress)
	t := mail.NewEmail(to, toAddress)
	textContent := ""
	// TODO: Parse body to ensure links are valid
	message := mail.NewSingleEmail(f, subject, t, textContent, body)
	client := sendgrid.NewSendClient(config.Var.SENDGRID_API_KEY)
	res, err := client.Send(message)

	if res.StatusCode >= 400 {
		return common.StringError(errors.New(res.Body))
	}

	if err != nil {
		return common.StringError(err)
	}
	return nil
}
