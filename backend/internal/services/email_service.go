package services

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/ngwebwork/payflow/internal/config"
)

type EmailService struct{ host, port, username, password, from string }

func NewEmailService(cfg config.Config) *EmailService {
	return &EmailService{host: cfg.SMTPHost, port: cfg.SMTPPort, username: cfg.SMTPUsername, password: cfg.SMTPPassword, from: cfg.SMTPFrom}
}
func (s *EmailService) SendVerification(email, code string) error {
	return s.send(email, "Verify your PayFlow email", fmt.Sprintf("Your PayFlow verification code is: %s\n\nThis code expires in 10 minutes.", code))
}
func (s *EmailService) SendPasswordReset(email, token string) error {
	return s.send(email, "Reset your PayFlow password", fmt.Sprintf("Use this PayFlow password reset token: %s\n\nThis token expires in 15 minutes.", token))
}
func (s *EmailService) SendInvoice(email, invoiceNumber string, total float64, currency string, dueDate string) error {
	return s.send(email, "Invoice "+invoiceNumber+" from PayFlow", fmt.Sprintf("Your invoice %s is ready. Total due: %.2f %s\nDue date: %s", invoiceNumber, total, currency, dueDate))
}
func (s *EmailService) send(to, subject, body string) error {
	if s.host == "" || s.username == "" || s.password == "" {
		return fmt.Errorf("email delivery is not configured; set SMTP_HOST, SMTP_USERNAME, SMTP_PASSWORD, and SMTP_FROM")
	}
	message := []byte("From: " + s.from + "\r\nTo: " + to + "\r\nSubject: " + subject + "\r\n\r\n" + body + "\r\n")
	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	if s.port == "465" {
		return s.sendTLS(to, message, auth)
	}
	if err := smtp.SendMail(s.host+":"+s.port, auth, s.from, []string{to}, message); err != nil {
		return fmt.Errorf("send email: %w", err)
	}
	return nil
}

// sendTLS delivers over implicit TLS (port 465); smtp.SendMail only supports STARTTLS.
func (s *EmailService) sendTLS(to string, message []byte, auth smtp.Auth) error {
	conn, err := tls.Dial("tcp", s.host+":"+s.port, &tls.Config{ServerName: s.host})
	if err != nil {
		return fmt.Errorf("send email: %w", err)
	}
	defer conn.Close()
	client, err := smtp.NewClient(conn, s.host)
	if err != nil {
		return fmt.Errorf("send email: %w", err)
	}
	defer client.Close()
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("send email: %w", err)
	}
	if err := client.Mail(s.from); err != nil {
		return fmt.Errorf("send email: %w", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("send email: %w", err)
	}
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("send email: %w", err)
	}
	if _, err := writer.Write(message); err != nil {
		return fmt.Errorf("send email: %w", err)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("send email: %w", err)
	}
	return client.Quit()
}
