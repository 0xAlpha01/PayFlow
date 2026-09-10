package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Notification struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"userId" json:"-"`
	Type      string             `bson:"type" json:"type"`
	Title     string             `bson:"title" json:"title"`
	Message   string             `bson:"message" json:"message"`
	Read      bool               `bson:"read" json:"read"`
	Metadata  bson.M             `bson:"metadata,omitempty" json:"metadata,omitempty"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}
