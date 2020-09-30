package email

type PasswordTemplate struct {
	Token string `json:"token"`
}

type EmailSender interface {
	Send(to string, templateId interface{}) error
}
