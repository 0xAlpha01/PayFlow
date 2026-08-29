package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ngwebwork/payflow/internal/config"
)

type PaystackService struct {
	secret, baseURL string
	client          *http.Client
}
type PaystackInitializeResponse struct {
	AuthorizationURL string `json:"authorization_url"`
	Reference        string `json:"reference"`
}

func NewPaystackService(cfg config.Config) *PaystackService {
	return &PaystackService{secret: cfg.PaystackSecretKey, baseURL: strings.TrimRight(cfg.PaystackBaseURL, "/"), client: &http.Client{}}
}
func (s *PaystackService) Initialize(ctx context.Context, email, reference string, amountKobo int64) (PaystackInitializeResponse, error) {
	var result struct {
		Status  bool                       `json:"status"`
		Message string                     `json:"message"`
		Data    PaystackInitializeResponse `json:"data"`
	}
	payload, _ := json.Marshal(map[string]any{"email": email, "amount": amountKobo, "reference": reference, "currency": "NGN"})
	if err := s.request(ctx, http.MethodPost, "/transaction/initialize", payload, &result); err != nil {
		return PaystackInitializeResponse{}, err
	}
	if !result.Status {
		return PaystackInitializeResponse{}, fmt.Errorf("paystack: %s", safePaystackMessage(result.Message))
	}
	return result.Data, nil
}
func (s *PaystackService) Verify(ctx context.Context, reference string) (string, error) {
	var result struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    struct {
			Status string `json:"status"`
		} `json:"data"`
	}
	if err := s.request(ctx, http.MethodGet, "/transaction/verify/"+url.PathEscape(reference), nil, &result); err != nil {
		return "", err
	}
	if !result.Status {
		return "", fmt.Errorf("paystack: %s", safePaystackMessage(result.Message))
	}
	return result.Data.Status, nil
}
func (s *PaystackService) request(ctx context.Context, method, path string, body []byte, target any) error {
	if s.secret == "" {
		return fmt.Errorf("paystack is not configured")
	}
	req, err := http.NewRequestWithContext(ctx, method, s.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.secret)
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("paystack request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("paystack returned %s", resp.Status)
	}
	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("invalid paystack response: %w", err)
	}
	return nil
}
func safePaystackMessage(message string) string {
	if strings.TrimSpace(message) == "" {
		return "request was rejected"
	}
	return message
}
