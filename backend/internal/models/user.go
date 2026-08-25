package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FirstName     string             `bson:"firstName" json:"firstName"`
	LastName      string             `bson:"lastName" json:"lastName"`
	Email         string             `bson:"email" json:"email"`
	Phone         string             `bson:"phone,omitempty" json:"phone,omitempty"`
	Password      string             `bson:"password" json:"-"`
	BusinessName  string             `bson:"businessName,omitempty" json:"businessName,omitempty"`
	EmailVerified bool               `bson:"emailVerified" json:"emailVerified"`
	PhoneVerified bool               `bson:"phoneVerified" json:"phoneVerified"`
	IsActive      bool               `bson:"isActive" json:"isActive"`
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updatedAt" json:"updatedAt"`
}
