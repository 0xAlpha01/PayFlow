package services

import "log"

type EmailService struct{}

func NewEmailService() *EmailService { return &EmailService{} }
func (s *EmailService) SendVerification(email, code string) error {
	log.Printf("email verification queued for %s (provider integration pending)", email)
	return nil
}
func (s *EmailService) SendPasswordReset(email, token string) error {
	log.Printf("password reset queued for %s (provider integration pending)", email)
	return nil
}
