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

// initialTestBalance is credited exactly once, server-side, when a new account is
// created (VerifyEmail is the only code path that inserts a user document) — the
// frontend never supplies or influences this value.
const initialTestBalance = 50000

type AuthService struct {
	users         *repositories.UserRepository
	redis         *redis.Client
	jwt           *JWTService
	email         *EmailService
	phone         *PhoneService
	notifications *NotificationService
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
type PhoneVerification struct {
	Phone     string    `json:"phone"`
	CodeHash  string    `json:"codeHash"`
	Attempts  int       `json:"attempts"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewAuthService(users *repositories.UserRepository, r *redis.Client, jwt *JWTService, email *EmailService, phone *PhoneService, notifications *NotificationService) *AuthService {
	return &AuthService{users: users, redis: r, jwt: jwt, email: email, phone: phone, notifications: notifications}
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
		phone, err := utils.NormalizeNigerianPhone(input.Phone)
		if err != nil {
			return nil, err
		}
		input.Phone = phone
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
func (s *AuthService) VerifyEmail(ctx context.Context, email, verification string) (*models.User, string, string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	key := "pending_signup:" + email
	pendingData, err := s.redis.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, "", "", errors.New("verification code has expired. Please request a new code")
	}
	if err != nil {
		return nil, "", "", errors.New("unable to verify email")
	}
	var pending PendingSignup
	if err := json.Unmarshal([]byte(pendingData), &pending); err != nil {
		return nil, "", "", errors.New("invalid verification request")
	}
	if utils.CheckPassword(pending.VerificationHash, verification) != nil {
		return nil, "", "", errors.New("invalid verification code")
	}
	existing, err := s.users.ByEmail(ctx, pending.Email)
	if err != nil {
		return nil, "", "", err
	}
	if existing != nil {
		return nil, "", "", errors.New("email is already registered")
	}
	if pending.Phone != "" {
		existing, err = s.users.ByPhone(ctx, pending.Phone)
		if err != nil {
			return nil, "", "", err
		}
		if existing != nil {
			return nil, "", "", errors.New("phone number is already registered")
		}
	}
	user := &models.User{ID: primitive.NewObjectID(), FirstName: pending.FirstName, LastName: pending.LastName, Email: pending.Email, Phone: pending.Phone, Password: pending.PasswordHash, Balance: initialTestBalance, BusinessName: pending.BusinessName, EmailVerified: true, IsActive: true, CreatedAt: pending.CreatedAt, UpdatedAt: time.Now().UTC()}
	if err := s.users.Create(ctx, user); err != nil {
		return nil, "", "", err
	}
	if s.notifications != nil {
		_ = s.notifications.Create(ctx, user.ID, "account_created", "Welcome to PayFlow!", "Your account is ready and funded with a NGN 50,000.00 test balance.", nil)
	}
	if err := s.redis.Del(ctx, key).Err(); err != nil {
		return nil, "", "", err
	}
	access, err := s.jwt.Access(user.ID, user.Email)
	if err != nil {
		return nil, "", "", err
	}
	refresh, err := s.jwt.Refresh(user.ID, user.Email)
	if err != nil {
		return nil, "", "", err
	}
	claims, err := s.jwt.Parse(refresh, "refresh")
	if err != nil {
		return nil, "", "", err
	}
	if err := s.redis.Set(ctx, "session:"+claims.ID, user.ID.Hex(), 7*24*time.Hour).Err(); err != nil {
		return nil, "", "", err
	}
	return user, access, refresh, nil
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
func (s *AuthService) SendPhoneCode(ctx context.Context, userID primitive.ObjectID, requestedPhone string, resend bool) error {
	user, err := s.users.ByID(ctx, userID)
	if err != nil || user == nil {
		return errors.New("user not found")
	}
	if user.PhoneVerified {
		return errors.New("phone is already verified")
	}
	if user.Phone == "" {
		return errors.New("phone number is required")
	}
	phone, err := utils.NormalizeNigerianPhone(requestedPhone)
	if err != nil || phone != user.Phone {
		return errors.New("invalid phone number")
	}
	if !resend {
		if exists, err := s.redis.Exists(ctx, "phone_verification:"+userID.Hex()).Result(); err == nil && exists > 0 {
			return nil
		}
	}
	allowed, err := s.redis.SetNX(ctx, "phone_verification_resend:"+userID.Hex(), "1", time.Minute).Result()
	if err != nil {
		return err
	}
	if !allowed {
		return errors.New("please wait before requesting another code")
	}
	value := code()
	if value == "" {
		return errors.New("unable to generate verification code")
	}
	hash, err := utils.HashPassword(value)
	if err != nil {
		return err
	}
	data, err := json.Marshal(PhoneVerification{Phone: user.Phone, CodeHash: hash, CreatedAt: time.Now().UTC()})
	if err != nil {
		return err
	}
	if err := s.redis.Set(ctx, "phone_verification:"+userID.Hex(), data, 10*time.Minute).Err(); err != nil {
		return err
	}
	if err := s.phone.SendVerificationCode(user.Phone, value); err != nil {
		_ = s.redis.Del(ctx, "phone_verification:"+userID.Hex()).Err()
		return err
	}
	return nil
}

func (s *AuthService) VerifyPhone(ctx context.Context, userID primitive.ObjectID, value string) error {
	if len(value) != 6 {
		return errors.New("invalid verification code")
	}
	user, err := s.users.ByID(ctx, userID)
	if err != nil || user == nil {
		return errors.New("user not found")
	}
	if user.PhoneVerified {
		return errors.New("phone is already verified")
	}
	key := "phone_verification:" + userID.Hex()
	data, err := s.redis.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return errors.New("verification code has expired. Please request a new code")
	}
	if err != nil {
		return errors.New("unable to verify phone")
	}
	var verification PhoneVerification
	if err := json.Unmarshal([]byte(data), &verification); err != nil {
		return errors.New("invalid phone verification request")
	}
	if verification.Attempts >= 5 {
		return errors.New("too many attempts. Please request a new code")
	}
	if utils.CheckPassword(verification.CodeHash, value) != nil {
		verification.Attempts++
		updated, _ := json.Marshal(verification)
		ttl, ttlErr := s.redis.TTL(ctx, key).Result()
		if ttlErr == nil && ttl > 0 {
			_ = s.redis.Set(ctx, key, updated, ttl).Err()
		}
		if verification.Attempts >= 5 {
			return errors.New("too many attempts. Please request a new code")
		}
		return errors.New("invalid verification code")
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
