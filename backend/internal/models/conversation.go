package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Conversation struct {
	ID             primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	ParticipantIDs []primitive.ObjectID `bson:"participantIds" json:"participantIds"`
	PairKey        string               `bson:"pairKey" json:"-"`
	LastMessage    string               `bson:"lastMessage,omitempty" json:"lastMessage,omitempty"`
	LastMessageAt  time.Time            `bson:"lastMessageAt,omitempty" json:"lastMessageAt,omitempty"`
	LastSenderID   primitive.ObjectID   `bson:"lastSenderId,omitempty" json:"lastSenderId,omitempty"`
	CreatedAt      time.Time            `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time            `bson:"updatedAt" json:"updatedAt"`
}

// PublicUser is the safe subset of a User exposed to other users through chat/wallet
// search results — never the full User document (balance, verification flags, etc.).
type PublicUser struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	FirstName    string             `bson:"firstName" json:"firstName"`
	LastName     string             `bson:"lastName" json:"lastName"`
	Email        string             `bson:"email" json:"email"`
	BusinessName string             `bson:"businessName,omitempty" json:"businessName,omitempty"`
}

type ChatMessage struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ConversationID primitive.ObjectID `bson:"conversationId" json:"conversationId"`
	SenderID       primitive.ObjectID `bson:"senderId" json:"senderId"`
	ReceiverID     primitive.ObjectID `bson:"receiverId" json:"receiverId"`
	Content        string             `bson:"content" json:"content"`
	Read           bool               `bson:"read" json:"read"`
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
}
