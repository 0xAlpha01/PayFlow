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
