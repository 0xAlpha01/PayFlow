package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ngwebwork/payflow/internal/config"
	"github.com/ngwebwork/payflow/internal/database"
	"github.com/ngwebwork/payflow/internal/models"
	"github.com/ngwebwork/payflow/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	client, err := database.NewMongo(ctx, cfg.MongoURI)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())
	db := client.Database(cfg.MongoDatabase)
	hash, err := utils.HashPassword("TestPass123!")
	if err != nil {
		log.Fatal(err)
	}
	users := db.Collection("users")
	now := time.Now().UTC()
	for _, u := range []models.User{
		{ID: primitive.NewObjectID(), FirstName: "UI", LastName: "Tester", Email: "ui.tester@payflow.local", Phone: "+2348022220001", Password: hash, Balance: 42000, BusinessName: "UI Test Co", EmailVerified: true, PhoneVerified: true, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{ID: primitive.NewObjectID(), FirstName: "Nina", LastName: "Reyes", Email: "nina.reyes@payflow.local", Phone: "+2348022220002", Password: hash, Balance: 15000, EmailVerified: true, PhoneVerified: true, IsActive: true, CreatedAt: now, UpdatedAt: now},
	} {
		_, _ = users.DeleteOne(ctx, bson.M{"email": u.Email})
		if _, err := users.InsertOne(ctx, u); err != nil {
			log.Fatal(err)
		}
		fmt.Println("seeded", u.Email)
	}
}
