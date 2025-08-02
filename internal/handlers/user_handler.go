package handlers

import (
	"BecomeOverMan/internal/models"
	"BecomeOverMan/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if err := h.service.Register(user.Username, user.Email, user.PasswordHash); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func (h *UserHandler) Login(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	userID, err := h.service.Login(user.Username, user.PasswordHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user_id": userID,
	})
}

// RegisterQuestRoutes sets up the routes for quest handling with Gin
func RegisterUserRoutes(router *gin.RouterGroup, userService *services.UserService) {
	handler := NewUserHandler(userService)

	userGroup := router.Group("/user")
	{
		userGroup.POST("/login", handler.Login)
		userGroup.POST("/register", handler.Register)
	}
}
