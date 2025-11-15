package send_email

import (
	"auth/configs"
	dto "auth/internal/app/domain/core/dto/services/send-email"
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
)

func SendEmail(email dto.EmailDTO) error {
	emailBody, err := parseTemplate(email.Template, email.Data)
	if err != nil {
		return fmt.Errorf("failed to parse and execute template: %w", err)
	}
	email.Data.SenderMail = configs.Config.Mail.FromAddress

	message := prepareEmailMessage(email, emailBody)

	if err = sendEmailSMTP(email.Data.SenderMail, email.Email, message); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	fmt.Println("Email sent successfully to:", email.Email)

	return nil
}

func parseTemplate(templateName string, data interface{}) (string, error) {
	tmplPath := fmt.Sprintf("templates/%s.html", templateName)
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return "", fmt.Errorf("failed to parse template file: %w", err)
	}

	var body bytes.Buffer
	if err = tmpl.Execute(&body, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return body.String(), nil
}

func prepareEmailMessage(email dto.EmailDTO, body string) []byte {
	return []byte(fmt.Sprintf(
		"To: %s\r\nSubject: %s\r\nFrom: %s\r\nMIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n%s",
		email.Email, email.Data.Subject, email.Data.SenderMail, body,
	))
}

func sendEmailSMTP(from, to string, msg []byte) error {
	auth := smtp.PlainAuth("", configs.Config.Mail.Username, configs.Config.Mail.Password, configs.Config.Mail.Host)
	smtpAddr := fmt.Sprintf("%s:%d", configs.Config.Mail.Host, configs.Config.Mail.Port)

	return smtp.SendMail(smtpAddr, auth, from, []string{to}, msg)
}
