package services

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ngwebwork/payflow/internal/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type JWTService struct{ cfg config.Config }
type Claims struct {
	UserID    string `json:"userId"`
	Email     string `json:"email"`
	TokenType string `json:"tokenType"`
	jwt.RegisteredClaims
}

func NewJWTService(cfg config.Config) *JWTService { return &JWTService{cfg: cfg} }
func (s *JWTService) create(userID primitive.ObjectID, email, tokenType string, expiry time.Duration, secret string) (string, error) {
	now := time.Now()
	claims := Claims{UserID: userID.Hex(), Email: email, TokenType: tokenType, RegisteredClaims: jwt.RegisteredClaims{IssuedAt: jwt.NewNumericDate(now), ExpiresAt: jwt.NewNumericDate(now.Add(expiry)), ID: primitive.NewObjectID().Hex()}}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}
func (s *JWTService) Access(userID primitive.ObjectID, email string) (string, error) {
	return s.create(userID, email, "access", s.cfg.AccessExpiry, s.cfg.AccessSecret)
}
func (s *JWTService) Refresh(userID primitive.ObjectID, email string) (string, error) {
	return s.create(userID, email, "refresh", s.cfg.RefreshExpiry, s.cfg.RefreshSecret)
}
func (s *JWTService) Parse(token, expectedType string) (*Claims, error) {
	secret := s.cfg.AccessSecret
	if expectedType == "refresh" {
		secret = s.cfg.RefreshSecret
	}
	parsed, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil || !parsed.Valid {
		return nil, errors.New("invalid or expired token")
	}
	claims, ok := parsed.Claims.(*Claims)
	if !ok || claims.TokenType != expectedType {
		return nil, errors.New("invalid token type")
	}
	return claims, nil
}
