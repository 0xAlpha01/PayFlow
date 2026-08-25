package services

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/ngwebwork/payflow/internal/config"
)

type EmailService struct{ host, port, username, password, from string }

func NewEmailService(cfg config.Config) *EmailService {
	return &EmailService{host: cfg.SMTPHost, port: cfg.SMTPPort, username: cfg.SMTPUsername, password: cfg.SMTPPassword, from: cfg.SMTPFrom}
}
func (s *EmailService) SendVerification(email, code string) error {
	return s.send(email, "Verify your PayFlow email", fmt.Sprintf("Your PayFlow verification code is: %s\n\nThis code expires in 10 minutes.", code), code)
}
func (s *EmailService) SendPasswordReset(email, token string) error {
	return s.send(email, "Reset your PayFlow password", fmt.Sprintf("Use this PayFlow password reset token: %s\n\nThis token expires in 15 minutes.", token), "")
}
func (s *EmailService) send(to, subject, body, developmentCode string) error {
	if s.host == "" || s.username == "" || s.password == "" {
		if developmentCode != "" {
			log.Printf("development email code for %s: %s", to, developmentCode)
		} else {
			log.Printf("development password reset email queued for %s", to)
		}
		return nil
	}
	message := []byte("From: " + s.from + "\r\nTo: " + to + "\r\nSubject: " + subject + "\r\n\r\n" + body + "\r\n")
	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	if err := smtp.SendMail(s.host+":"+s.port, auth, s.from, []string{to}, message); err != nil {
		return fmt.Errorf("send email: %w", err)
	}
	return nil
}
