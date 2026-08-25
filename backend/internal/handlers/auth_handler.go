package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/ngwebwork/payflow/internal/services"
	"github.com/ngwebwork/payflow/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strings"
)

type AuthHandler struct{ auth *services.AuthService }

func NewAuthHandler(auth *services.AuthService) *AuthHandler { return &AuthHandler{auth: auth} }
func (h *AuthHandler) Signup(c *gin.Context) {
	var input services.SignupInput
	if c.ShouldBindJSON(&input) != nil {
		utils.Failure(c, 400, "invalid request body", nil)
		return
	}
	user, err := h.auth.Signup(c.Request.Context(), input)
	if err != nil {
		code := http.StatusBadRequest
		if strings.Contains(err.Error(), "already") {
			code = http.StatusConflict
		}
		utils.Failure(c, code, err.Error(), nil)
		return
	}
	utils.Success(c, 201, "Account created successfully. Please verify your email.", gin.H{"user": user})
}
func (h *AuthHandler) Login(c *gin.Context) {
	var input services.LoginInput
	if c.ShouldBindJSON(&input) != nil {
		utils.Failure(c, 400, "invalid request body", nil)
		return
	}
	user, access, refresh, err := h.auth.Login(c.Request.Context(), input)
	if err != nil {
		utils.Failure(c, 401, err.Error(), nil)
		return
	}
	utils.Success(c, 200, "Login successful", gin.H{"user": user, "accessToken": access, "refreshToken": refresh})
}
func (h *AuthHandler) Refresh(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refreshToken"`
	}
	if c.ShouldBindJSON(&body) != nil || body.RefreshToken == "" {
		utils.Failure(c, 400, "refreshToken is required", nil)
		return
	}
	access, err := h.auth.Refresh(c.Request.Context(), body.RefreshToken)
	if err != nil {
		utils.Failure(c, 401, err.Error(), nil)
		return
	}
	utils.Success(c, 200, "Token refreshed", gin.H{"accessToken": access})
}
func (h *AuthHandler) Logout(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refreshToken"`
	}
	if c.ShouldBindJSON(&body) != nil || body.RefreshToken == "" {
		utils.Failure(c, 400, "refreshToken is required", nil)
		return
	}
	if err := h.auth.Logout(c.Request.Context(), body.RefreshToken); err != nil {
		utils.Failure(c, 401, "invalid refresh token", nil)
		return
	}
	utils.Success(c, 200, "Logout successful", nil)
}
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var body struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	if c.ShouldBindJSON(&body) != nil {
		utils.Failure(c, 400, "email and code are required", nil)
		return
	}
	if err := h.auth.VerifyEmail(c, body.Email, body.Code); err != nil {
		utils.Failure(c, 400, err.Error(), nil)
		return
	}
	utils.Success(c, 200, "Email verified", nil)
}
func (h *AuthHandler) ResendEmail(c *gin.Context) {
	var body struct {
		Email string `json:"email"`
	}
	if c.ShouldBindJSON(&body) != nil {
		utils.Failure(c, 400, "email is required", nil)
		return
	}
	if err := h.auth.ResendEmail(c, body.Email); err != nil {
		utils.Failure(c, 500, "unable to resend verification", nil)
		return
	}
	utils.Success(c, 200, "Verification email sent", nil)
}
func (h *AuthHandler) SendPhone(c *gin.Context) {
	id, ok := userID(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}
	if err := h.auth.SendPhone(c, id); err != nil {
		utils.Failure(c, 400, err.Error(), nil)
		return
	}
	utils.Success(c, 200, "Phone verification code sent", nil)
}
func (h *AuthHandler) VerifyPhone(c *gin.Context) {
	id, ok := userID(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}
	var body struct {
		Code string `json:"code"`
	}
	if c.ShouldBindJSON(&body) != nil {
		utils.Failure(c, 400, "code is required", nil)
		return
	}
	if err := h.auth.VerifyPhone(c, id, body.Code); err != nil {
		utils.Failure(c, 400, err.Error(), nil)
		return
	}
	utils.Success(c, 200, "Phone verified", nil)
}
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var body struct {
		Email string `json:"email"`
	}
	if c.ShouldBindJSON(&body) != nil {
		utils.Failure(c, 400, "email is required", nil)
		return
	}
	_ = h.auth.ForgotPassword(c, body.Email)
	utils.Success(c, 200, "If that email exists, reset instructions have been sent", nil)
}
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var body struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}
	if c.ShouldBindJSON(&body) != nil {
		utils.Failure(c, 400, "token and password are required", nil)
		return
	}
	if err := h.auth.ResetPassword(c, body.Token, body.Password); err != nil {
		utils.Failure(c, 400, err.Error(), nil)
		return
	}
	utils.Success(c, 200, "Password reset successful", nil)
}
func (h *AuthHandler) Me(c *gin.Context) {
	id, ok := userID(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}
	user, err := h.auth.User(c, id)
	if err != nil || user == nil {
		utils.Failure(c, 404, "user not found", nil)
		return
	}
	utils.Success(c, 200, "Current user", gin.H{"user": user})
}
func userID(c *gin.Context) (primitive.ObjectID, bool) {
	value, ok := c.Get("userID")
	if !ok {
		return primitive.NilObjectID, false
	}
	id, err := primitive.ObjectIDFromHex(value.(string))
	return id, err == nil
}
