package mailer

import (
	"bytes"
	"embed"
	"log"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/go-mail/mail/v2"
)

type Mailer struct {
	dialer *mail.Dialer
	sender string
}

//go:embed "templates"
var templateFS embed.FS

var (
	host     = os.Getenv("SMTP_HOST")
	port     = os.Getenv("SMTP_PORT")
	username = os.Getenv("SMTP_USERNAME")
	password = os.Getenv("SMTP_PASSWORD")
	sender   = os.Getenv("SMTP_SENDER")
)

func New() Mailer {
	portInt, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalf("Invalid port number: %v", err)
	}
	dialer := mail.NewDialer(host, portInt, username, password)
	dialer.Timeout = 5 * time.Second

	return Mailer{dialer: dialer, sender: sender}
}

func (m Mailer) Send(recipient, templateFile string, data any) error {
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	msg := mail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	for i := 1; i <= 3; i++ {
		err = m.dialer.DialAndSend(msg)
		// If it worked, return nil
		if nil == err {
			return nil
		}

		// If it didn't work, sleep for a short time and retry
		time.Sleep(500 * time.Millisecond)
	}

	return err
}
