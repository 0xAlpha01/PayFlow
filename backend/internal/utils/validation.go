package utils

import (
	"errors"
	"net/mail"
	"strings"
)

func Required(values ...string) error {
	for _, value := range values {
		if strings.TrimSpace(value) == "" {
			return errors.New("required fields are missing")
		}
	}
	return nil
}
func Email(value string) error {
	parsed, err := mail.ParseAddress(value)
	if err != nil || parsed.Address != value {
		return errors.New("invalid email address")
	}
	return nil
}

// NormalizeNigerianPhone converts common Nigerian mobile formats to E.164.
func NormalizeNigerianPhone(value string) (string, error) {
	phone := strings.NewReplacer(" ", "", "-", "", "(", "", ")", "").Replace(strings.TrimSpace(value))
	switch {
	case strings.HasPrefix(phone, "+234"):
		phone = phone[1:]
	case strings.HasPrefix(phone, "0"):
		phone = "234" + phone[1:]
	}
	if len(phone) != 13 || !strings.HasPrefix(phone, "234") || phone[3] != '8' {
		return "", errors.New("invalid Nigerian phone number")
	}
	for _, char := range phone {
		if char < '0' || char > '9' {
			return "", errors.New("invalid Nigerian phone number")
		}
	}
	return "+" + phone, nil
}
