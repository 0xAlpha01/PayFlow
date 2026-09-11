package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port, MongoURI, MongoDatabase, RedisAddr, RedisPassword, ClientURL string
	SMTPHost, SMTPPort, SMTPUsername, SMTPPassword, SMTPFrom           string
	SMSProvider, SMSAPIKey, SMSSenderID                                string
	PaystackSecretKey, PaystackPublicKey, PaystackBaseURL               string
	RedisDB                                                            int
	AccessSecret, RefreshSecret                                        string
	AccessExpiry, RefreshExpiry                                        time.Duration
}

func Load() (Config, error) {
	// Load the .env file from the backend directory when present. In production
	// (e.g. Render, Railway) there is no .env file — env vars are set directly
	// on the platform, so a missing file here is expected, not an error.
	if err := godotenv.Load(".env"); err != nil && !os.IsNotExist(err) {
		return Config{}, fmt.Errorf("failed to load .env file: %w", err)
	}

	// SMTP configuration debug.
	// The actual password is NEVER printed.
	fmt.Println("=== SMTP DEBUG ===")
	fmt.Println("SMTP_HOST:", os.Getenv("SMTP_HOST"))
	fmt.Println("SMTP_PORT:", os.Getenv("SMTP_PORT"))
	fmt.Println("SMTP_USERNAME:", os.Getenv("SMTP_USERNAME"))
	fmt.Println("SMTP_PASSWORD SET:", os.Getenv("SMTP_PASSWORD") != "")
	fmt.Println("SMTP_FROM:", os.Getenv("SMTP_FROM"))
	fmt.Println("==================")

	access, err := time.ParseDuration(
		get("ACCESS_TOKEN_EXPIRY", "15m"),
	)
	if err != nil {
		return Config{}, fmt.Errorf(
			"invalid access expiry: %w",
			err,
		)
	}

	refresh, err := time.ParseDuration(
		get("REFRESH_TOKEN_EXPIRY", "168h"),
	)
	if err != nil {
		return Config{}, fmt.Errorf(
			"invalid refresh expiry: %w",
			err,
		)
	}

	cfg := Config{
		Port:          get("PORT", "8080"),
		MongoURI:      get("MONGODB_URI", "mongodb://localhost:27017"),
		MongoDatabase: get("MONGODB_DATABASE", "payflow"),

		RedisAddr:     get("REDIS_ADDR", "localhost:6379"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),

		ClientURL: get(
			"CLIENT_URL",
			"http://localhost:5173",
		),

		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     get("SMTP_PORT", "587"),
		SMTPUsername: os.Getenv("SMTP_USERNAME"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		SMTPFrom: get(
			"SMTP_FROM",
			os.Getenv("SMTP_USERNAME"),
		),
		SMSProvider: os.Getenv("SMS_PROVIDER"),
		SMSAPIKey:   os.Getenv("SMS_API_KEY"),
		SMSSenderID: os.Getenv("SMS_SENDER_ID"),
		PaystackSecretKey: os.Getenv("PAYSTACK_SECRET_KEY"),
		PaystackPublicKey: os.Getenv("PAYSTACK_PUBLIC_KEY"),
		PaystackBaseURL: get("PAYSTACK_BASE_URL", "https://api.paystack.co"),

		AccessSecret:  os.Getenv("JWT_ACCESS_SECRET"),
		RefreshSecret: os.Getenv("JWT_REFRESH_SECRET"),

		AccessExpiry:  access,
		RefreshExpiry: refresh,
	}

	if cfg.AccessSecret == "" || cfg.RefreshSecret == "" {
		return Config{}, fmt.Errorf(
			"JWT_ACCESS_SECRET and JWT_REFRESH_SECRET must be configured",
		)
	}

	return cfg, nil
}

func get(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
