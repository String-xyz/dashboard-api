package emailer

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"text/template"

	"github.com/String-xyz/dashboard-api/config"
	libcommon "github.com/String-xyz/go-lib/v2/common"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Emailer interface {
	SendInviteEmail(ctx context.Context, email string, token string, inviteId string, userName string) error
	SendPasswordResetEmail(ctx context.Context, email string, token string, userName string) error
}

//go:embed templates/*
var templatesFS embed.FS

type emailer struct {
}

func New() Emailer {
	return &emailer{}
}

func (e emailer) SendPasswordResetEmail(ctx context.Context, email string, token string, userName string) error {
	link := config.Var.BASE_DASHBOARD_URL + "/members/password-reset/" + token // TODO: Change URL to Password Reset page and include resetToken

	tmpl, err := template.ParseFS(templatesFS, "templates/password_reset.tpl")
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "password_reset.tpl", map[string]interface{}{
		"link":     link,
		"userName": userName,
	})
	if err != nil {
		return err
	}

	from := mail.NewEmail("String API", config.Var.AUTH_EMAIL_ADDRESS)
	subject := "String API Password Reset"
	to := mail.NewEmail(userName, email)
	return sendEmail(ctx, from, subject, to, "", buf.String())
}

func (e emailer) SendInviteEmail(ctx context.Context, email string, token string, inviteId string, userName string) error {
	link := config.Var.BASE_DASHBOARD_URL + "/invite/" + inviteId + "?token=" + token

	tmpl, err := template.ParseFS(templatesFS, "templates/invite.tpl")
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "invite.tpl", map[string]interface{}{
		"link":     link,
		"userName": userName,
	})
	if err != nil {
		return err
	}

	from := mail.NewEmail("String API", config.Var.AUTH_EMAIL_ADDRESS)
	subject := "New String API User"
	to := mail.NewEmail(userName, email)
	return sendEmail(ctx, from, subject, to, "", buf.String())
}

func sendEmail(ctx context.Context, from *mail.Email, subject string, to *mail.Email, text string, html string) error {
	message := mail.NewSingleEmail(from, subject, to, text, html)
	client := sendgrid.NewSendClient(config.Var.SENDGRID_API_KEY)
	res, err := client.Send(message)
	if res.StatusCode >= 400 {
		return libcommon.StringError(errors.New(res.Body))
	}
	if err != nil {
		return libcommon.StringError(err)
	}
	return nil
}
