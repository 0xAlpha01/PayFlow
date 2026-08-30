package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ngwebwork/payflow/internal/config"
)

// PhoneService keeps SMS-provider details out of the authentication service.
type PhoneService struct {
	provider string
	apiKey   string
	senderID string
	client   *http.Client
}

func NewPhoneService(cfg config.Config) *PhoneService {
	return &PhoneService{provider: strings.ToLower(cfg.SMSProvider), apiKey: cfg.SMSAPIKey, senderID: cfg.SMSSenderID, client: &http.Client{Timeout: 10 * time.Second}}
}

func (s *PhoneService) SendVerificationCode(phone, code string) error {
	if s.provider != "termii" || s.apiKey == "" || s.senderID == "" {
		return errors.New("SMS delivery is not configured")
	}
	payload, err := json.Marshal(map[string]string{
		"to": strings.TrimPrefix(phone, "+"), "from": s.senderID, "sms": fmt.Sprintf("Your PayFlow verification code is %s. It expires in 10 minutes.", code), "type": "plain", "channel": "generic", "api_key": s.apiKey,
	})
	if err != nil {
		return err
	}
	request, err := http.NewRequest(http.MethodPost, "https://api.ng.termii.com/api/sms/send", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := s.client.Do(request)
	if err != nil {
		return errors.New("unable to deliver SMS")
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return errors.New("SMS provider failed to send verification code")
	}
	return nil
}
