package services

import (
	"context"
	"errors"
	"math"
	"strings"
	"time"

	"github.com/ngwebwork/payflow/internal/models"
	"github.com/ngwebwork/payflow/internal/repositories"
	"github.com/ngwebwork/payflow/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SendTransferInput struct {
	RecipientEmail  string  `json:"recipientEmail"`
	Amount          float64 `json:"amount"`
	Note            string  `json:"note"`
	ClientReference string  `json:"clientReference"`
}
type TransferService struct {
	transfers     *repositories.TransferRepository
	users         *repositories.UserRepository
	notifications *NotificationService
}

func NewTransferService(transfers *repositories.TransferRepository, users *repositories.UserRepository, notifications *NotificationService) *TransferService {
	return &TransferService{transfers: transfers, users: users, notifications: notifications}
}
func (s *TransferService) Balance(ctx context.Context, userID primitive.ObjectID) (float64, error) {
	return s.transfers.Balance(ctx, userID)
}
func name(u *models.User) string {
	name := strings.TrimSpace(u.FirstName + " " + u.LastName)
	if name == "" {
		return u.Email
	}
	return name
}
func (s *TransferService) Send(ctx context.Context, senderID primitive.ObjectID, input SendTransferInput) (*models.Transfer, error) {
	input.RecipientEmail = strings.ToLower(strings.TrimSpace(input.RecipientEmail))
	input.ClientReference = strings.TrimSpace(input.ClientReference)
	if err := utils.Required(input.RecipientEmail); err != nil {
		return nil, err
	}
	if err := utils.Email(input.RecipientEmail); err != nil {
		return nil, err
	}
	if input.Amount <= 0 || math.IsNaN(input.Amount) || math.IsInf(input.Amount, 0) {
		return nil, errors.New("amount must be greater than zero")
	}
	if existing, err := s.transfers.ByClientReference(ctx, senderID, input.ClientReference); err != nil {
		return nil, err
	} else if existing != nil {
		return existing, nil
	}
	sender, err := s.users.ByID(ctx, senderID)
	if err != nil || sender == nil {
		return nil, errors.New("sender account not found")
	}
	receiver, err := s.users.ByEmail(ctx, input.RecipientEmail)
	if err != nil {
		return nil, err
	}
	if receiver == nil {
		return nil, errors.New("no PayFlow user found with that email")
	}
	if receiver.ID == senderID {
		return nil, errors.New("you cannot send money to yourself")
	}
	transfer := &models.Transfer{ID: primitive.NewObjectID(), SenderID: senderID, ReceiverID: receiver.ID, Amount: input.Amount, Currency: "NGN", Status: "completed", Note: strings.TrimSpace(input.Note), ClientReference: input.ClientReference, CreatedAt: time.Now().UTC()}
	if err := s.transfers.Execute(ctx, transfer); err != nil {
		if errors.Is(err, repositories.ErrInsufficientBalance) {
			_ = s.notifications.Create(ctx, senderID, "transfer_failed", "Transfer failed", "Your transfer of "+money(input.Amount)+" to "+name(receiver)+" failed: insufficient balance.", bson.M{"amount": input.Amount, "recipientEmail": receiver.Email})
			return nil, errors.New("insufficient balance")
		}
		return nil, err
	}
	_ = s.notifications.Create(ctx, senderID, "transfer_sent", "Transfer sent", "You sent "+money(input.Amount)+" to "+name(receiver)+".", bson.M{"transferId": transfer.ID, "amount": input.Amount})
	_ = s.notifications.Create(ctx, receiver.ID, "money_received", "Money received", "You received "+money(input.Amount)+" from "+name(sender)+".", bson.M{"transferId": transfer.ID, "amount": input.Amount})
	return transfer, nil
}
func (s *TransferService) List(ctx context.Context, userID primitive.ObjectID, page, limit int) (repositories.TransferPage, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return s.transfers.List(ctx, userID, page, limit)
}
func money(amount float64) string {
	return utils.FormatNaira(amount)
}
