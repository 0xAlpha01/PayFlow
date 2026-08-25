package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"math/big"

	"github.com/ngwebwork/payflow/internal/models"
	"github.com/ngwebwork/payflow/internal/repositories"
	"github.com/ngwebwork/payflow/internal/utils"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthService struct {
	users *repositories.UserRepository
	redis *redis.Client
	jwt   *JWTService
	email *EmailService
	phone *PhoneService
}
type SignupInput struct {
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Password     string `json:"password"`
	BusinessName string `json:"businessName"`
}
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type PendingSignup struct {
	FirstName        string    `json:"firstName"`
	LastName         string    `json:"lastName"`
	BusinessName     string    `json:"businessName"`
	Email            string    `json:"email"`
	Phone            string    `json:"phone"`
	PasswordHash     string    `json:"passwordHash"`
	VerificationHash string    `json:"verificationHash"`
	CreatedAt        time.Time `json:"createdAt"`
}

func NewAuthService(users *repositories.UserRepository, r *redis.Client, jwt *JWTService, email *EmailService, phone *PhoneService) *AuthService {
	return &AuthService{users: users, redis: r, jwt: jwt, email: email, phone: phone}
}
func code() string {
	number, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%06d", number.Int64())
}
func token() string { b := make([]byte, 32); _, _ = rand.Read(b); return hex.EncodeToString(b) }
func (s *AuthService) Signup(ctx context.Context, input SignupInput) (*models.User, error) {
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	if err := utils.Required(input.FirstName, input.LastName, input.Email, input.Password); err != nil {
		return nil, err
	}
	if err := utils.Email(input.Email); err != nil {
		return nil, err
	}
	if err := utils.ValidatePassword(input.Password); err != nil {
		return nil, err
	}
	existing, err := s.users.ByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email is already registered")
	}
	if input.Phone != "" {
		existing, err = s.users.ByPhone(ctx, input.Phone)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, errors.New("phone number is already registered")
		}
	}
	hash, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}
	verification := code()
	if verification == "" {
		return nil, errors.New("unable to generate verification code")
	}
	verificationHash, err := utils.HashPassword(verification)
	if err != nil {
		return nil, err
	}
	pending := PendingSignup{FirstName: input.FirstName, LastName: input.LastName, BusinessName: input.BusinessName, Email: input.Email, Phone: input.Phone, PasswordHash: hash, VerificationHash: verificationHash, CreatedAt: time.Now().UTC()}
	pendingData, err := json.Marshal(pending)
	if err != nil {
		return nil, err
	}
	key := "pending_signup:" + input.Email
	if err := s.redis.Set(ctx, key, pendingData, 10*time.Minute).Err(); err != nil {
		return nil, err
	}
	if err := s.email.SendVerification(input.Email, verification); err != nil {
		_ = s.redis.Del(ctx, key).Err()
		return nil, err
	}
	return nil, nil
}
func (s *AuthService) Login(ctx context.Context, input LoginInput) (*models.User, string, string, error) {
	user, err := s.users.ByEmail(ctx, strings.ToLower(strings.TrimSpace(input.Email)))
	if err != nil || user == nil || utils.CheckPassword(user.Password, input.Password) != nil {
		return nil, "", "", errors.New("invalid email or password")
	}
	if !user.IsActive {
		return nil, "", "", errors.New("account is inactive")
	}
	access, err := s.jwt.Access(user.ID, user.Email)
	if err != nil {
		return nil, "", "", err
	}
	refresh, err := s.jwt.Refresh(user.ID, user.Email)
	if err != nil {
		return nil, "", "", err
	}
	claims, _ := s.jwt.Parse(refresh, "refresh")
	if err := s.redis.Set(ctx, "session:"+claims.ID, user.ID.Hex(), 7*24*time.Hour).Err(); err != nil {
		return nil, "", "", err
	}
	return user, access, refresh, nil
}
func (s *AuthService) Refresh(ctx context.Context, refresh string) (string, error) {
	claims, err := s.jwt.Parse(refresh, "refresh")
	if err != nil {
		return "", err
	}
	if s.redis.Get(ctx, "session:"+claims.ID).Err() != nil {
		return "", errors.New("session has been revoked")
	}
	id, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		return "", errors.New("invalid user id")
	}
	return s.jwt.Access(id, claims.Email)
}
func (s *AuthService) Logout(ctx context.Context, refresh string) error {
	claims, err := s.jwt.Parse(refresh, "refresh")
	if err != nil {
		return err
	}
	return s.redis.Del(ctx, "session:"+claims.ID).Err()
}
func (s *AuthService) VerifyEmail(ctx context.Context, email, verification string) error {
	email = strings.ToLower(strings.TrimSpace(email))
	key := "pending_signup:" + email
	pendingData, err := s.redis.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return errors.New("verification code has expired. Please request a new code")
	}
	if err != nil {
		return errors.New("unable to verify email")
	}
	var pending PendingSignup
	if err := json.Unmarshal([]byte(pendingData), &pending); err != nil {
		return errors.New("invalid verification request")
	}
	if utils.CheckPassword(pending.VerificationHash, verification) != nil {
		return errors.New("invalid verification code")
	}
	existing, err := s.users.ByEmail(ctx, pending.Email)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("email is already registered")
	}
	if pending.Phone != "" {
		existing, err = s.users.ByPhone(ctx, pending.Phone)
		if err != nil {
			return err
		}
		if existing != nil {
			return errors.New("phone number is already registered")
		}
	}
	user := &models.User{ID: primitive.NewObjectID(), FirstName: pending.FirstName, LastName: pending.LastName, Email: pending.Email, Phone: pending.Phone, Password: pending.PasswordHash, BusinessName: pending.BusinessName, EmailVerified: true, IsActive: true, CreatedAt: pending.CreatedAt, UpdatedAt: time.Now().UTC()}
	if err := s.users.Create(ctx, user); err != nil {
		return err
	}
	return s.redis.Del(ctx, key).Err()
}
func (s *AuthService) ResendEmail(ctx context.Context, email string) error {
	email = strings.ToLower(strings.TrimSpace(email))
	key := "pending_signup:" + email
	data, err := s.redis.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return errors.New("verification request has expired")
	}
	if err != nil {
		return err
	}
	allowed, err := s.redis.SetNX(ctx, "pending_signup_resend:"+email, "1", time.Minute).Result()
	if err != nil {
		return err
	}
	if !allowed {
		return errors.New("verification email recently sent; please wait before requesting another")
	}
	var pending PendingSignup
	if err := json.Unmarshal([]byte(data), &pending); err != nil {
		return errors.New("invalid verification request")
	}
	value := code()
	valueHash, err := utils.HashPassword(value)
	if err != nil {
		return err
	}
	pending.VerificationHash = valueHash
	pending.CreatedAt = time.Now().UTC()
	updated, err := json.Marshal(pending)
	if err != nil {
		return err
	}
	if err := s.redis.Set(ctx, key, updated, 10*time.Minute).Err(); err != nil {
		return err
	}
	return s.email.SendVerification(pending.Email, value)
}
func (s *AuthService) SendPhone(ctx context.Context, phone string) error {
	user, err := s.users.ByPhone(ctx, phone)
	if err != nil || user == nil || user.Phone == "" {
		return errors.New("phone number is required")
	}
	value := code()
	if err := s.redis.Set(ctx, "verify:phone:"+phone, value, 10*time.Minute).Err(); err != nil {
		return err
	}
	return s.phone.SendCode(user.Phone, value)
}
func (s *AuthService) VerifyPhone(ctx context.Context, phone, value string) error {
	key := "verify:phone:" + phone
	if s.redis.Get(ctx, key).Val() != value {
		return errors.New("invalid or expired verification code")
	}
	user, err := s.users.ByPhone(ctx, phone)
	if err != nil || user == nil {
		return errors.New("invalid phone verification request")
	}
	if err := s.users.Update(ctx, user.ID, bson.M{"phoneVerified": true}); err != nil {
		return err
	}
	return s.redis.Del(ctx, key).Err()
}
func (s *AuthService) ForgotPassword(ctx context.Context, email string) error {
	user, err := s.users.ByEmail(ctx, email)
	if err != nil || user == nil {
		return nil
	}
	value := token()
	if err := s.redis.Set(ctx, "reset:"+value, user.ID.Hex(), 15*time.Minute).Err(); err != nil {
		return err
	}
	return s.email.SendPasswordReset(user.Email, value)
}
func (s *AuthService) ResetPassword(ctx context.Context, value, password string) error {
	if err := utils.ValidatePassword(password); err != nil {
		return err
	}
	idValue, err := s.redis.Get(ctx, "reset:"+value).Result()
	if err != nil {
		return errors.New("invalid or expired reset token")
	}
	id, err := primitive.ObjectIDFromHex(idValue)
	if err != nil {
		return errors.New("invalid reset token")
	}
	hash, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	if err := s.users.Update(ctx, id, bson.M{"password": hash}); err != nil {
		return err
	}
	return s.redis.Del(ctx, "reset:"+value).Err()
}

func (s *AuthService) User(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	return s.users.ByID(ctx, id)
}
