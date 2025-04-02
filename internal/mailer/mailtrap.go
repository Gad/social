package mailer



import (
	"fmt"
	gomail "gopkg.in/mail.v2"
)

type mailtrapMailer struct {
	fromEmail string
	apiKey string
	
}


func NewMailtrap(apiKey, fromEmail string) *mailtrapMailer{
	return &mailtrapMailer{
		fromEmail: fromEmail,
		apiKey: apiKey,
	}
}

func (m *mailtrapMailer) Send(templateFile string, username string, email string, data any, isSandbox bool) (int,error) {
    
	// template parsing and building

	
	
	
	// Create a new message
    message := gomail.NewMessage()

    // Set email headers
    message.SetHeader("From", m.fromEmail)
    message.SetHeader("To", email)
    message.SetHeader("Subject", subject)

    // Set the plain-text version of the email
    //message.SetBody("text/plain", body)

    // Set the HTML version of the email
    message.AddAlternative("text/html", body)

    // Set up the SMTP dialer
    dialer := gomail.NewDialer("live.smtp.mailtrap.io", 587, "api", m.apiKey)

    // Send the email
    if err := dialer.DialAndSend(message); err != nil {
        return -1, err
    } 
	return 200, nil
}