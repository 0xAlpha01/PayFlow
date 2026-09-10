package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/ngwebwork/payflow/internal/middleware"
	"github.com/ngwebwork/payflow/internal/realtime"
	"github.com/ngwebwork/payflow/internal/services"
	"github.com/ngwebwork/payflow/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatHandler struct {
	service   *services.ChatService
	jwt       *services.JWTService
	hub       *realtime.Hub
	upgrader  websocket.Upgrader
	clientURL string
}

func NewChatHandler(service *services.ChatService, jwt *services.JWTService, hub *realtime.Hub, clientURL string) *ChatHandler {
	allowed := middleware.AllowedOrigins(clientURL)
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			if origin == "" {
				return true
			}
			for _, allow := range allowed {
				if origin == allow {
					return true
				}
			}
			return false
		},
	}
	return &ChatHandler{service: service, jwt: jwt, hub: hub, upgrader: upgrader, clientURL: clientURL}
}
func (h *ChatHandler) SearchUsers(c *gin.Context) {
	userID, ok := notificationUser(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}
	users, err := h.service.SearchUsers(c, userID, c.Query("q"))
	if err != nil {
		utils.Failure(c, 500, "unable to search users", nil)
		return
	}
	utils.Success(c, 200, "Users loaded", users)
}
func (h *ChatHandler) StartConversation(c *gin.Context) {
	userID, ok := notificationUser(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}
	var body struct {
		UserID string `json:"userId"`
	}
	if c.ShouldBindJSON(&body) != nil {
		utils.Failure(c, 400, "invalid request body", nil)
		return
	}
	peerID, err := primitive.ObjectIDFromHex(body.UserID)
	if err != nil {
		utils.Failure(c, 400, "invalid user id", nil)
		return
	}
	conversation, err := h.service.StartConversation(c, userID, peerID)
	if err != nil {
		utils.Failure(c, 400, err.Error(), nil)
		return
	}
	utils.Success(c, 201, "Conversation ready", conversation)
}
func (h *ChatHandler) ListConversations(c *gin.Context) {
	userID, ok := notificationUser(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}
	conversations, err := h.service.ListConversations(c, userID)
	if err != nil {
		utils.Failure(c, 500, "unable to load conversations", nil)
		return
	}
	utils.Success(c, 200, "Conversations loaded", conversations)
}
func chatConversationID(c *gin.Context) (primitive.ObjectID, bool) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	return id, err == nil
}
func (h *ChatHandler) Messages(c *gin.Context) {
	userID, ok := notificationUser(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}
	conversationID, ok := chatConversationID(c)
	if !ok {
		utils.Failure(c, 400, "invalid conversation id", nil)
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	result, err := h.service.Messages(c, conversationID, userID, page, limit)
	if err != nil {
		utils.Failure(c, chatErrorStatus(err), err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Messages loaded", "data": result.Items, "pagination": gin.H{"page": result.Page, "limit": result.Limit, "total": result.Total, "totalPages": result.TotalPages}})
}
func (h *ChatHandler) SendMessage(c *gin.Context) {
	userID, ok := notificationUser(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}
	conversationID, ok := chatConversationID(c)
	if !ok {
		utils.Failure(c, 400, "invalid conversation id", nil)
		return
	}
	var body struct {
		Content string `json:"content"`
	}
	if c.ShouldBindJSON(&body) != nil {
		utils.Failure(c, 400, "invalid request body", nil)
		return
	}
	message, err := h.service.SendMessage(c, conversationID, userID, body.Content)
	if err != nil {
		utils.Failure(c, chatErrorStatus(err), err.Error(), nil)
		return
	}
	utils.Success(c, 201, "Message sent", message)
}
func (h *ChatHandler) MarkRead(c *gin.Context) {
	userID, ok := notificationUser(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}
	conversationID, ok := chatConversationID(c)
	if !ok {
		utils.Failure(c, 400, "invalid conversation id", nil)
		return
	}
	if err := h.service.MarkRead(c, conversationID, userID); err != nil {
		utils.Failure(c, chatErrorStatus(err), err.Error(), nil)
		return
	}
	utils.Success(c, 200, "Conversation marked as read", nil)
}
func (h *ChatHandler) UnreadCount(c *gin.Context) {
	userID, ok := notificationUser(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}
	count, err := h.service.UnreadCount(c, userID)
	if err != nil {
		utils.Failure(c, 500, "unable to load unread count", nil)
		return
	}
	utils.Success(c, 200, "Unread count loaded", gin.H{"count": count})
}
func chatErrorStatus(err error) int {
	if errors.Is(err, services.ErrNotParticipant) {
		return http.StatusForbidden
	}
	if strings.Contains(err.Error(), "not found") {
		return http.StatusNotFound
	}
	return http.StatusBadRequest
}

// Socket upgrades an authenticated connection to a websocket and registers it in the
// hub so notifications and chat messages can be pushed to this user in real time.
// Browsers cannot set custom headers on the WebSocket handshake, so the access token
// travels as a query parameter here instead of the Authorization header.
func (h *ChatHandler) Socket(c *gin.Context) {
	token := c.Query("token")
	claims, err := h.jwt.Parse(token, "access")
	if err != nil {
		utils.Failure(c, 401, "invalid or expired access token", nil)
		return
	}
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	userID := claims.UserID
	h.hub.Register(userID, conn)
	defer func() {
		h.hub.Unregister(userID, conn)
		_ = conn.Close()
	}()
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			return
		}
	}
}
