package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID                        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FirstName                 string             `bson:"firstName" json:"firstName"`
	LastName                  string             `bson:"lastName" json:"lastName"`
	Email                     string             `bson:"email" json:"email"`
	Phone                     string             `bson:"phone,omitempty" json:"phone,omitempty"`
	Password                  string             `bson:"password" json:"-"`
	BusinessName              string             `bson:"businessName,omitempty" json:"businessName,omitempty"`
	BusinessEmail             string             `bson:"businessEmail,omitempty" json:"businessEmail,omitempty"`
	BusinessPhone             string             `bson:"businessPhone,omitempty" json:"businessPhone,omitempty"`
	BusinessAddress           string             `bson:"businessAddress,omitempty" json:"businessAddress,omitempty"`
	BusinessDescription       string             `bson:"businessDescription,omitempty" json:"businessDescription,omitempty"`
	BusinessWebsite           string             `bson:"businessWebsite,omitempty" json:"businessWebsite,omitempty"`
	EmailNotifications        *bool              `bson:"emailNotifications,omitempty" json:"-"`
	PaymentNotifications      *bool              `bson:"paymentNotifications,omitempty" json:"-"`
	InvoiceNotifications      *bool              `bson:"invoiceNotifications,omitempty" json:"-"`
	SubscriptionNotifications *bool              `bson:"subscriptionNotifications,omitempty" json:"-"`
	SecurityNotifications     *bool              `bson:"securityNotifications,omitempty" json:"-"`
	EmailVerified             bool               `bson:"emailVerified" json:"emailVerified"`
	PhoneVerified             bool               `bson:"phoneVerified" json:"phoneVerified"`
	IsActive                  bool               `bson:"isActive" json:"isActive"`
	CreatedAt                 time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt                 time.Time          `bson:"updatedAt" json:"updatedAt"`
}
