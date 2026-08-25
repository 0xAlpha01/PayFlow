package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port, MongoURI, MongoDatabase, RedisAddr, RedisPassword, ClientURL string
	RedisDB                                                            int
	AccessSecret, RefreshSecret                                        string
	AccessExpiry, RefreshExpiry                                        time.Duration
}

func Load() (Config, error) {
	_ = godotenv.Load()
	access, err := time.ParseDuration(get("ACCESS_TOKEN_EXPIRY", "15m"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid access expiry: %w", err)
	}
	refresh, err := time.ParseDuration(get("REFRESH_TOKEN_EXPIRY", "168h"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid refresh expiry: %w", err)
	}
	cfg := Config{Port: get("PORT", "8080"), MongoURI: get("MONGODB_URI", "mongodb://localhost:27017"), MongoDatabase: get("MONGODB_DATABASE", "payflow"), RedisAddr: get("REDIS_ADDR", "localhost:6379"), RedisPassword: os.Getenv("REDIS_PASSWORD"), ClientURL: get("CLIENT_URL", "http://localhost:5173"), AccessSecret: os.Getenv("JWT_ACCESS_SECRET"), RefreshSecret: os.Getenv("JWT_REFRESH_SECRET"), AccessExpiry: access, RefreshExpiry: refresh}
	if cfg.AccessSecret == "" || cfg.RefreshSecret == "" {
		return Config{}, fmt.Errorf("JWT_ACCESS_SECRET and JWT_REFRESH_SECRET must be configured")
	}
	return cfg, nil
}
func get(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
