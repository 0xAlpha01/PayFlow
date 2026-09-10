package services

import (
	"context"
	"errors"

	"github.com/ngwebwork/payflow/internal/repositories"
	"github.com/ngwebwork/payflow/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SettingsService struct{ users *repositories.UserRepository }

func NewSettingsService(users *repositories.UserRepository) *SettingsService {
	return &SettingsService{users: users}
}
func (s *SettingsService) Profile(ctx context.Context, id primitive.ObjectID) (bson.M, error) {
	user, err := s.users.ByID(ctx, id)
	if err != nil || user == nil {
		return nil, err
	}
	return bson.M{"firstName": user.FirstName, "lastName": user.LastName, "email": user.Email, "phone": user.Phone, "emailVerified": user.EmailVerified, "phoneVerified": user.PhoneVerified}, nil
}
func (s *SettingsService) UpdateProfile(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	safe := bson.M{}
	for _, key := range []string{"firstName", "lastName", "phone"} {
		if value, ok := update[key]; ok {
			safe[key] = value
		}
	}
	return s.users.Update(ctx, id, safe)
}
func (s *SettingsService) Business(ctx context.Context, id primitive.ObjectID) (bson.M, error) {
	user, err := s.users.ByID(ctx, id)
	if err != nil || user == nil {
		return nil, err
	}
	return bson.M{"businessName": user.BusinessName, "businessEmail": user.BusinessEmail, "businessPhone": user.BusinessPhone, "businessAddress": user.BusinessAddress, "businessDescription": user.BusinessDescription, "businessWebsite": user.BusinessWebsite}, nil
}
func (s *SettingsService) UpdateBusiness(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	safe := bson.M{}
	for _, key := range []string{"businessName", "businessEmail", "businessPhone", "businessAddress", "businessDescription", "businessWebsite"} {
		if value, ok := update[key]; ok {
			safe[key] = value
		}
	}
	return s.users.Update(ctx, id, safe)
}
func boolOr(value *bool, fallback bool) bool {
	if value == nil {
		return fallback
	}
	return *value
}
func (s *SettingsService) Notifications(ctx context.Context, id primitive.ObjectID) (bson.M, error) {
	user, err := s.users.ByID(ctx, id)
	if err != nil || user == nil {
		return nil, err
	}
	return bson.M{"emailNotifications": boolOr(user.EmailNotifications, true), "paymentNotifications": boolOr(user.PaymentNotifications, true), "invoiceNotifications": boolOr(user.InvoiceNotifications, true), "subscriptionNotifications": boolOr(user.SubscriptionNotifications, true), "securityNotifications": boolOr(user.SecurityNotifications, true)}, nil
}
func (s *SettingsService) UpdateNotifications(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	safe := bson.M{}
	for _, key := range []string{"emailNotifications", "paymentNotifications", "invoiceNotifications", "subscriptionNotifications", "securityNotifications"} {
		if value, ok := update[key].(bool); ok {
			safe[key] = value
		}
	}
	return s.users.Update(ctx, id, safe)
}
func (s *SettingsService) ChangePassword(ctx context.Context, id primitive.ObjectID, current, next string) error {
	user, err := s.users.ByID(ctx, id)
	if err != nil || user == nil {
		return errors.New("account not found")
	}
	if utils.CheckPassword(user.Password, current) != nil {
		return errors.New("current password is incorrect")
	}
	if err := utils.ValidatePassword(next); err != nil {
		return err
	}
	hash, err := utils.HashPassword(next)
	if err != nil {
		return err
	}
	return s.users.Update(ctx, id, bson.M{"password": hash})
}
