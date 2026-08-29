package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type InvoiceItem struct {
	Description string  `bson:"description" json:"description"`
	Quantity    int     `bson:"quantity" json:"quantity"`
	Price       float64 `bson:"price" json:"price"`
}
type Invoice struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID        primitive.ObjectID `bson:"userId" json:"-"`
	CustomerID    primitive.ObjectID `bson:"customerId,omitempty" json:"customerId,omitempty"`
	InvoiceNumber string             `bson:"invoiceNumber" json:"invoiceNumber"`
	Items         []InvoiceItem      `bson:"items" json:"items"`
	Subtotal      float64            `bson:"subtotal" json:"subtotal"`
	Tax           float64            `bson:"tax" json:"tax"`
	Discount      float64            `bson:"discount" json:"discount"`
	Total         float64            `bson:"total" json:"total"`
	Currency      string             `bson:"currency" json:"currency"`
	Status        string             `bson:"status" json:"status"`
	DueDate       time.Time          `bson:"dueDate" json:"dueDate"`
	Notes         string             `bson:"notes,omitempty" json:"notes,omitempty"`
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updatedAt" json:"updatedAt"`
}
