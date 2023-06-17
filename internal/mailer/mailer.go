package mailer

import (
	"fmt"
	"time"

	"github.com/alpden550/go-ecommerce-stripe/internal/configs"
	mail "github.com/xhit/go-simple-mail/v2"
)

func serveSMTPClient(configer configs.BaseConfiger) (*mail.SMTPClient, error) {
	config := configer.GetConfig()
	errorLog := configer.GetErrorLog()

	mailer := mail.NewSMTPClient()
	mailer.Host = config.SMTP.Host
	mailer.Port = config.SMTP.Port
	mailer.Username = config.SMTP.Username
	mailer.Password = config.SMTP.Password
	mailer.Encryption = mail.EncryptionTLS
	mailer.KeepAlive = false
	mailer.ConnectTimeout = 10 * time.Second
	mailer.SendTimeout = 10 * time.Second

	smtpClient, err := mailer.Connect()
	if err != nil {
		errorLog.Printf("%w", fmt.Errorf("%e", err))
		return nil, err
	}

	return smtpClient, nil
}

func SendEmail(
	configer configs.BaseConfiger,
	from,
	to,
	subject,
	textTmpl, htmlTmpl string,
	attachments []string,
	data interface{},
) error {
	errorLog := configer.GetErrorLog()

	formattedMessage, err := renderTemplate(htmlTmpl, data)
	if err != nil {
		errorLog.Printf("%w", fmt.Errorf("%e", err))
		return err
	}

	plainMessage, err := renderTemplate(textTmpl, data)
	if err != nil {
		errorLog.Printf("%w", fmt.Errorf("%e", err))
		return err
	}

	smtp, err := serveSMTPClient(configer)
	if err != nil {
		errorLog.Printf("%w", fmt.Errorf("%e", err))
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(from).
		AddTo(to).
		SetSubject(subject)

	email.SetBody(mail.TextHTML, formattedMessage)
	email.AddAlternative(mail.TextPlain, plainMessage)

	if len(attachments) > 0 {
		for _, attachment := range attachments {
			email.AddAttachment(attachment)
		}
	}

	if err = email.Send(smtp); err != nil {
		errorLog.Printf("%w", fmt.Errorf("%e", err))
		return err
	}

	return nil
}
