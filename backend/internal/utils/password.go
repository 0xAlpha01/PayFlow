package utils

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"unicode"
)

func ValidatePassword(password string) error {
	var lower, upper, number, special bool
	if len([]rune(password)) < 8 {
		return errors.New("password must be at least 8 characters and include uppercase, lowercase, number, and special character")
	}
	for _, character := range password {
		switch {
		case unicode.IsLower(character):
			lower = true
		case unicode.IsUpper(character):
			upper = true
		case unicode.IsNumber(character):
			number = true
		default:
			special = true
		}
	}
	if !lower || !upper || !number || !special {
		return errors.New("password must be at least 8 characters and include uppercase, lowercase, number, and special character")
	}
	return nil
}
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}
func CheckPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
