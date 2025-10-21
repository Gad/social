package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"time"

	gomail "gopkg.in/mail.v2"
)

type mailtrapMailer struct {
	fromEmail    string
	apiKey       string
	smtpAddr     string
	smtpPort     int
	smtpUsername string
	maxRetries   int
}

func NewMailtrap(apiKey, fromEmail, smtpAddr, smtpUsername string, smtpPort, maxRetries int) *mailtrapMailer {
	return &mailtrapMailer{
		fromEmail:    fromEmail,
		apiKey:       apiKey,
		smtpAddr:     smtpAddr,
		smtpPort:     smtpPort,
		smtpUsername: smtpUsername,
		maxRetries:   maxRetries,
	}
}

func (m *mailtrapMailer) Send(templateFile string, username string, email string, data any, isSandbox bool) error {

	if isSandbox {
		return nil
	}
	// template parsing and building
	subject, body, err := m.htmlEmailFromTemplate(templateFile, data)
	if err != nil {
		return err
	}

	// Create a new message
	err = m.mailTrapDialing(email, subject, body)
	if err != nil {
		return err
	}
	return nil
}

func (m *mailtrapMailer) htmlEmailFromTemplate(templateFile string, data any) (*bytes.Buffer, *bytes.Buffer, error) {
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return nil, nil, err
	}
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return nil, nil, err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return nil, nil, err
	}

	return subject, body, nil
}

func (m *mailtrapMailer) mailTrapDialing(destinationEmail string, subject, body *bytes.Buffer) error {

	message := gomail.NewMessage()

	// Set email headers
	message.SetHeader("From", m.fromEmail)
	message.SetHeader("To", destinationEmail)
	message.SetHeader("Subject", subject.String())

	// Set the plain-text version of the email
	//message.SetBody("text/plain", body)

	// Set the HTML version of the email
	message.AddAlternative("text/html", body.String())

	// Set up the SMTP dialer
	dialer := gomail.NewDialer(m.smtpAddr, m.smtpPort, m.smtpUsername, m.apiKey)

	// Send the email, maxReties attempts

	var RetryError error

	for i := range m.maxRetries {
		RetryError = dialer.DialAndSend(message)
		if RetryError != nil {
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}
		return nil

	}
	return fmt.Errorf("failed to send verification email after %d attempts with error :%v", m.maxRetries, RetryError)

}
