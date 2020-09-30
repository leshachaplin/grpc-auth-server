package email

import (
	"fmt"
	"net/smtp"
)

type SMTPEmail struct {
	from string
	username string
	password string
	host string
}

func NewSMTPEmail(from, username, password, host string) *SMTPEmail {
	return &SMTPEmail{
		from:     from,
		username: username,
		password: password,
		host: host,
	}
}
//smtp.gmail.com

func(e *SMTPEmail) Send(to string, data interface{}) error {
	d, ok := data.(PasswordTemplate)
	if !ok{
		return fmt.Errorf("unable to convert data")
	}
	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", e.username, e.password, e.host),
		e.from, []string{to}, []byte(fmt.Sprintf("<a href='%s'>click here</a>",d.Token)))
	return err
}
