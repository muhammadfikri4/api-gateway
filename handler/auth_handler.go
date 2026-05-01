package handler

import (
	"net/http"
	"time"

	"api-gateway/config"
	"api-gateway/model"
	"api-gateway/repository"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepo *repository.UserRepository
	cfg      *config.Config
}

func NewAuthHandler(userRepo *repository.UserRepository, cfg *config.Config) *AuthHandler {
	return &AuthHandler{userRepo: userRepo, cfg: cfg}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Success: false, Message: err.Error()})
		return
	}

	user, err := h.userRepo.FindByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.Response{Success: false, Message: "Invalid email or password"})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusForbidden, model.Response{Success: false, Message: "Account is deactivated"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, model.Response{Success: false, Message: "Invalid email or password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(h.cfg.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Success: false, Message: "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Success: true,
		Message: "Login successful",
		Data: gin.H{
			"token": tokenString,
			"user":  user,
		},
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Success: false, Message: err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Success: false, Message: "Failed to hash password"})
		return
	}

	role := "user"
	if req.Role != "" {
		role = req.Role
	}

	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     role,
		IsActive: true,
	}

	if err := h.userRepo.Create(user); err != nil {
		c.JSON(http.StatusConflict, model.Response{Success: false, Message: "Email already exists"})
		return
	}

	c.JSON(http.StatusCreated, model.Response{
		Success: true,
		Message: "User registered successfully",
		Data:    user,
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, _ := c.Get("user_id")
	user, err := h.userRepo.FindByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, model.Response{Success: false, Message: "User not found"})
		return
	}
	c.JSON(http.StatusOK, model.Response{Success: true, Data: user})
}
