package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ngwebwork/payflow/internal/config"
	"github.com/ngwebwork/payflow/internal/database"
	"github.com/ngwebwork/payflow/internal/handlers"
	"github.com/ngwebwork/payflow/internal/repositories"
	"github.com/ngwebwork/payflow/internal/routes"
	"github.com/ngwebwork/payflow/internal/services"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	mongo, err := database.NewMongo(ctx, cfg.MongoURI)
	if err != nil {
		log.Fatalf("mongodb connection failed: %v", err)
	}
	defer func() { _ = mongo.Disconnect(context.Background()) }()
	redis, err := database.NewRedis(ctx, cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Fatalf("redis connection failed: %v", err)
	}
	defer redis.Close()
	repository, err := repositories.NewUserRepository(mongo.Database(cfg.MongoDatabase))
	if err != nil {
		log.Fatalf("user repository setup failed: %v", err)
	}
	jwt := services.NewJWTService(cfg)
	auth := services.NewAuthService(repository, redis, jwt, services.NewEmailService(cfg), services.NewPhoneService(cfg))
	handler := handlers.NewAuthHandler(auth)
	dashboardHandler := handlers.NewDashboardHandler(services.NewDashboardService(repositories.NewDashboardRepository(mongo.Database(cfg.MongoDatabase))))
	resources := handlers.NewResourceHandler(services.NewResourceService(repositories.NewResourceRepository(mongo.Database(cfg.MongoDatabase))))
	customerRepository, err := repositories.NewCustomerRepository(mongo.Database(cfg.MongoDatabase))
	if err != nil {
		log.Fatalf("customer repository setup failed: %v", err)
	}
	customers := handlers.NewCustomerHandler(services.NewCustomerService(customerRepository))
	settings := handlers.NewSettingsHandler(services.NewSettingsService(repository))
	paymentRepository, err := repositories.NewPaymentRepository(mongo.Database(cfg.MongoDatabase))
	if err != nil {
		log.Fatalf("payment repository setup failed: %v", err)
	}
	payments := handlers.NewPaymentHandler(services.NewPaymentService(paymentRepository, services.NewPaystackService(cfg)), cfg.PaystackSecretKey)
	invoiceRepository, err := repositories.NewInvoiceRepository(mongo.Database(cfg.MongoDatabase))
	if err != nil {
		log.Fatalf("invoice repository setup failed: %v", err)
	}
	invoices := handlers.NewInvoiceHandler(services.NewInvoiceService(invoiceRepository, services.NewEmailService(cfg)))
	subscriptionRepository, err := repositories.NewSubscriptionRepository(mongo.Database(cfg.MongoDatabase))
	if err != nil {
		log.Fatalf("subscription repository setup failed: %v", err)
	}
	subscriptions := handlers.NewSubscriptionHandler(services.NewSubscriptionService(subscriptionRepository))
	analytics := handlers.NewAnalyticsHandler(services.NewAnalyticsService(repositories.NewAnalyticsRepository(mongo.Database(cfg.MongoDatabase))))
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	routes.Setup(router, handler, dashboardHandler, resources, customers, settings, payments, invoices, subscriptions, analytics, jwt, cfg.ClientURL)
	server := &http.Server{Addr: ":" + cfg.Port, Handler: router}
	go func() {
		log.Printf("PayFlow API listening on :%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()
	<-ctx.Done()
	shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdown); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}
