package mailer

import (
	"fmt"
	"github.com/alpden550/go-ecommerce-stripe/internal/configs"
	mail "github.com/xhit/go-simple-mail/v2"
	"time"
)

func serveSMTPClient(api *configs.ApiApplication) (*mail.SMTPClient, error) {
	mailer := mail.NewSMTPClient()
	mailer.Host = api.Config.SMTP.Host
	mailer.Port = api.Config.SMTP.Port
	mailer.Username = api.Config.SMTP.Username
	mailer.Password = api.Config.SMTP.Password
	mailer.Encryption = mail.EncryptionTLS
	mailer.KeepAlive = false
	mailer.ConnectTimeout = 10 * time.Second
	mailer.SendTimeout = 10 * time.Second

	smtpClient, err := mailer.Connect()
	if err != nil {
		api.ErrorLog.Printf("%w", fmt.Errorf("%e", err))
		return nil, err
	}

	return smtpClient, nil
}

func SendEmail(api *configs.ApiApplication, from, to, subject, textTmpl, htmlTmpl string, data interface{}) error {
	formattedMessage, err := renderTemplate(htmlTmpl, data)
	if err != nil {
		api.ErrorLog.Printf("%w", fmt.Errorf("%e", err))
		return err
	}

	plainMessage, err := renderTemplate(textTmpl, data)
	if err != nil {
		api.ErrorLog.Printf("%w", fmt.Errorf("%e", err))
		return err
	}

	smtp, err := serveSMTPClient(api)
	if err != nil {
		api.ErrorLog.Printf("%w", fmt.Errorf("%e", err))
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(from).
		AddTo(to).
		SetSubject(subject)

	email.SetBody(mail.TextHTML, formattedMessage)
	email.AddAlternative(mail.TextPlain, plainMessage)

	if err = email.Send(smtp); err != nil {
		api.ErrorLog.Printf("%w", fmt.Errorf("%e", err))
		return err
	}

	return nil
}
