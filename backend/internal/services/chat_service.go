package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/ngwebwork/payflow/internal/models"
	"github.com/ngwebwork/payflow/internal/realtime"
	"github.com/ngwebwork/payflow/internal/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrNotParticipant = errors.New("you are not a participant in this conversation")

type ChatService struct {
	repo          *repositories.ChatRepository
	notifications *NotificationService
	hub           *realtime.Hub
}

func NewChatService(repo *repositories.ChatRepository, notifications *NotificationService, hub *realtime.Hub) *ChatService {
	return &ChatService{repo: repo, notifications: notifications, hub: hub}
}
func (s *ChatService) SearchUsers(ctx context.Context, currentUserID primitive.ObjectID, query string) ([]models.User, error) {
	return s.repo.SearchUsers(ctx, currentUserID, strings.TrimSpace(query), 20)
}
func (s *ChatService) StartConversation(ctx context.Context, userID, peerID primitive.ObjectID) (*models.Conversation, error) {
	if userID == peerID {
		return nil, errors.New("you cannot message yourself")
	}
	peer, err := s.repo.User(ctx, peerID)
	if err != nil {
		return nil, err
	}
	if peer == nil {
		return nil, errors.New("user not found")
	}
	return s.repo.GetOrCreateConversation(ctx, userID, peerID)
}
func (s *ChatService) ListConversations(ctx context.Context, userID primitive.ObjectID) ([]repositories.ConversationSummary, error) {
	return s.repo.ListConversations(ctx, userID)
}

// authorize loads a conversation and confirms userID is one of its participants,
// so a user can never read or post into a conversation they are not part of.
func (s *ChatService) authorize(ctx context.Context, conversationID, userID primitive.ObjectID) (*models.Conversation, error) {
	conversation, err := s.repo.ConversationByID(ctx, conversationID)
	if err != nil {
		return nil, err
	}
	if conversation == nil {
		return nil, errors.New("conversation not found")
	}
	for _, participant := range conversation.ParticipantIDs {
		if participant == userID {
			return conversation, nil
		}
	}
	return nil, ErrNotParticipant
}
func (s *ChatService) Messages(ctx context.Context, conversationID, userID primitive.ObjectID, page, limit int) (repositories.MessagePage, error) {
	if _, err := s.authorize(ctx, conversationID, userID); err != nil {
		return repositories.MessagePage{}, err
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 200 {
		limit = 50
	}
	return s.repo.ListMessages(ctx, conversationID, page, limit)
}
func (s *ChatService) SendMessage(ctx context.Context, conversationID, senderID primitive.ObjectID, content string) (*models.ChatMessage, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, errors.New("message cannot be empty")
	}
	if len(content) > 4000 {
		return nil, errors.New("message is too long")
	}
	conversation, err := s.authorize(ctx, conversationID, senderID)
	if err != nil {
		return nil, err
	}
	var receiverID primitive.ObjectID
	for _, participant := range conversation.ParticipantIDs {
		if participant != senderID {
			receiverID = participant
		}
	}
	message := &models.ChatMessage{ID: primitive.NewObjectID(), ConversationID: conversationID, SenderID: senderID, ReceiverID: receiverID, Content: content, Read: false, CreatedAt: time.Now().UTC()}
	if err := s.repo.SendMessage(ctx, message); err != nil {
		return nil, err
	}
	if s.hub != nil {
		s.hub.Send(receiverID.Hex(), map[string]any{"type": "message", "data": message})
	}
	sender, err := s.repo.User(ctx, senderID)
	if err == nil && sender != nil {
		preview := content
		if len(preview) > 80 {
			preview = preview[:80] + "…"
		}
		_ = s.notifications.Create(ctx, receiverID, "chat_message", "New message from "+name(sender), preview, bson.M{"conversationId": conversationID})
	}
	return message, nil
}
func (s *ChatService) MarkRead(ctx context.Context, conversationID, userID primitive.ObjectID) error {
	if _, err := s.authorize(ctx, conversationID, userID); err != nil {
		return err
	}
	return s.repo.MarkRead(ctx, conversationID, userID)
}
func (s *ChatService) UnreadCount(ctx context.Context, userID primitive.ObjectID) (int64, error) {
	return s.repo.UnreadCount(ctx, userID)
}
