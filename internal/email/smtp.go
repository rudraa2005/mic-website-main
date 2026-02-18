package email

import "net/smtp"

type SMTPService struct {
	from     string
	password string
	host     string
	port     string
}

func NewSMTPService(from, password, host, port string) *SMTPService {
	return &SMTPService{
		from:     from,
		password: password,
		host:     host,
		port:     port,
	}
}

func (s *SMTPService) Send(to string, subject string, body string) error {
	msg := "From: " + s.from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	auth := smtp.PlainAuth("", s.from, s.password, s.host)

	return smtp.SendMail(s.host+":"+s.port, auth, s.from, []string{to}, []byte(msg))
}
