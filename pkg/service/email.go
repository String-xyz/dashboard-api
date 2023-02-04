package service

import (
	"log"
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
	message := mail.NewSingleEmail(f, subject, t, textContent, body)
	log.Printf("\n\nemail message: %+v\n", message)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	log.Printf("\n\nemail client: %+v\n", client)
	resp, err := client.Send(message)
	if err != nil {
		return common.StringError(err)
	}
	log.Printf("\n\nemail resp: %+v\n", resp)
	return nil
}
