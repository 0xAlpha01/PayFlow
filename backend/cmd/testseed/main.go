package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ngwebwork/payflow/internal/config"
	"github.com/ngwebwork/payflow/internal/database"
	"github.com/ngwebwork/payflow/internal/models"
	"github.com/ngwebwork/payflow/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := database.NewMongo(ctx, cfg.MongoURI)
	if err != nil {
		panic(err)
	}
	db := client.Database(cfg.MongoDatabase)
	hash, _ := utils.HashPassword("TestPass123!")
	user := models.User{
		ID: primitive.NewObjectID(), FirstName: "Currency", LastName: "Checker",
		Email: "currency.checker@payflow.local", Phone: "+2348033330099",
		Password: hash, Balance: 50000, BusinessName: "Currency Check Co",
		EmailVerified: true, PhoneVerified: true, IsActive: true,
		CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(),
	}
	if _, err := db.Collection("users").InsertOne(ctx, user); err != nil {
		panic(err)
	}
	fmt.Println("SEEDED_OK", user.ID.Hex())
}
