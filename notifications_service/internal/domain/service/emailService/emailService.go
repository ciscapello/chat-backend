package emailservice

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/ciscapello/notification_service/application/config"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

type EmailConfig struct {
	From string
	Pass string
}

type EmailService struct {
	logger *zap.Logger
	config EmailConfig
	dialer *gomail.Dialer
}

func New(config config.Config, logger *zap.Logger) *EmailService {
	d := gomail.NewDialer("smtp.gmail.com", 587, config.EmailAddress, config.EmailPassword)

	return &EmailService{
		logger: logger,
		dialer: d,
		config: EmailConfig{
			From: config.EmailAddress,
			Pass: config.EmailPassword,
		},
	}
}

func (e *EmailService) SendCodeToUser(code string, email string) {
	e.logger.Info("Sending email")

	msg := gomail.NewMessage()

	msg.SetHeader("From", e.config.From)
	// msg.SetHeader("To", email)
	msg.SetHeader("To", "ciscapello@gmail.com")
	msg.SetHeader("Subject", "Code")
	msg.SetBody("text/plain", code)

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}

	data := struct {
		Code string
	}{
		Code: code,
	}

	// Построение пути к файлу шаблона
	templatePath := filepath.Join(wd, "internal", "domain", "service", "emailService", "code.html")

	htmlBody, err := template.ParseFiles(templatePath)
	if err != nil {
		e.logger.Error("failed to parse html body", zap.Error(err))
	}

	fmt.Println(htmlBody)

	var body bytes.Buffer
	if err := htmlBody.Execute(&body, data); err != nil {
		e.logger.Error("failed to execute html template", zap.Error(err))
	}

	fmt.Println(body)

	msg.SetBody("text/html", body.String())

	e.dialer.DialAndSend(msg)

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
	log.Println("Successfully sended to " + "ciscapello@gmail.com")

}
