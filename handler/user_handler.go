package handler

import (
	"net/http"
	"strconv"

	"api-gateway/model"
	"api-gateway/repository"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userRepo *repository.UserRepository
}

func NewUserHandler(userRepo *repository.UserRepository) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

func (h *UserHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	users, total, err := h.userRepo.FindAll(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Success: false, Message: "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, model.PaginatedResponse{
		Success: true,
		Message: "Users fetched successfully",
		Data:    users,
		Total:   total,
		Page:    page,
		Limit:   limit,
	})
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Success: false, Message: "Invalid ID"})
		return
	}

	user, err := h.userRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, model.Response{Success: false, Message: "User not found"})
		return
	}

	c.JSON(http.StatusOK, model.Response{Success: true, Data: user})
}

func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Success: false, Message: "Invalid ID"})
		return
	}

	user, err := h.userRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, model.Response{Success: false, Message: "User not found"})
		return
	}

	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Success: false, Message: err.Error()})
		return
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := h.userRepo.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Success: false, Message: "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, model.Response{Success: true, Message: "User updated successfully", Data: user})
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Success: false, Message: "Invalid ID"})
		return
	}

	if err := h.userRepo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Success: false, Message: "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, model.Response{Success: true, Message: "User deleted successfully"})
}

func (h *UserHandler) ResetPassword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Success: false, Message: "Invalid ID"})
		return
	}

	user, err := h.userRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, model.Response{Success: false, Message: "User not found"})
		return
	}

	var req struct {
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Success: false, Message: err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Success: false, Message: "Failed to hash password"})
		return
	}

	user.Password = string(hashedPassword)
	if err := h.userRepo.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Success: false, Message: "Failed to reset password"})
		return
	}

	c.JSON(http.StatusOK, model.Response{Success: true, Message: "Password reset successfully"})
}
