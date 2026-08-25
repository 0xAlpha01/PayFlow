package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Subscription struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID          primitive.ObjectID `bson:"userId" json:"-"`
	CustomerID      primitive.ObjectID `bson:"customerId,omitempty" json:"customerId,omitempty"`
	Plan            string             `bson:"plan" json:"plan"`
	Amount          float64            `bson:"amount" json:"amount"`
	Currency        string             `bson:"currency" json:"currency"`
	BillingCycle    string             `bson:"billingCycle" json:"billingCycle"`
	Status          string             `bson:"status" json:"status"`
	NextPaymentDate time.Time          `bson:"nextPaymentDate" json:"nextPaymentDate"`
	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt"`
}
