package services

import (
	"context"
	"time"

	"github.com/ngwebwork/payflow/internal/models"
	"github.com/ngwebwork/payflow/internal/realtime"
	"github.com/ngwebwork/payflow/internal/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationService struct {
	repo *repositories.NotificationRepository
	hub  *realtime.Hub
}

func NewNotificationService(repo *repositories.NotificationRepository, hub *realtime.Hub) *NotificationService {
	return &NotificationService{repo: repo, hub: hub}
}

// Create persists a notification and, if the user has a live connection, pushes it
// immediately so the bell updates without the client polling.
func (s *NotificationService) Create(ctx context.Context, userID primitive.ObjectID, kind, title, message string, metadata bson.M) error {
	n := &models.Notification{ID: primitive.NewObjectID(), UserID: userID, Type: kind, Title: title, Message: message, Metadata: metadata, CreatedAt: time.Now().UTC()}
	if err := s.repo.Create(ctx, n); err != nil {
		return err
	}
	if s.hub != nil {
		s.hub.Send(userID.Hex(), map[string]any{"type": "notification", "data": n})
	}
	return nil
}
func (s *NotificationService) List(ctx context.Context, userID primitive.ObjectID, page, limit int) (repositories.NotificationPage, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return s.repo.List(ctx, userID, page, limit)
}
func (s *NotificationService) UnreadCount(ctx context.Context, userID primitive.ObjectID) (int64, error) {
	return s.repo.UnreadCount(ctx, userID)
}
func (s *NotificationService) MarkRead(ctx context.Context, userID, id primitive.ObjectID) error {
	return s.repo.MarkRead(ctx, userID, id)
}
func (s *NotificationService) MarkAllRead(ctx context.Context, userID primitive.ObjectID) error {
	return s.repo.MarkAllRead(ctx, userID)
}
