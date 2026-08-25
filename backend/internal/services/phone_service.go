package services

import "log"

type PhoneService struct{}

func NewPhoneService() *PhoneService { return &PhoneService{} }
func (s *PhoneService) SendCode(phone, code string) error {
	log.Printf("phone verification queued for %s (provider integration pending)", phone)
	return nil
}
