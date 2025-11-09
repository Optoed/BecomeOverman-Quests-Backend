package handlers

import (
	"BecomeOverMan/internal/models"
	"BecomeOverMan/internal/services"
	"BecomeOverMan/pkg/middleware"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	service *services.UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user by username, email, and password
// @Tags user
// @Accept json
// @Produce json
// @Param user body models.RegisterRequest true "User credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if len(req.Username) == 0 || len(req.Password) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and Password can't be empty"})
		return
	}

	if err := h.service.Register(req.Username, req.Email, req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// Login godoc
// @Summary Login a user
// @Description Login with username and password, returns JWT
// @Tags user
// @Accept json
// @Produce json
// @Param user body models.LoginRequest true "User credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var request models.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userID, err := h.service.Login(request.Username, request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := services.GenerateJWT(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user_id": userID,
		"token":   token,
	})
}

// AddFriend godoc
// @Summary Add friend
// @Description Add another user as friend
// @Tags friends
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param friend_id path int true "Friend ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /friends/{friend_id} [post]
func (h *UserHandler) AddFriend(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	friendID, err := strconv.Atoi(c.Param("friend_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid friend ID"})
		return
	}

	if err := h.service.AddFriend(userID, friendID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend added successfully"})
}

// GetFriends godoc
// @Summary Get user's friends
// @Description Get list of user's friends
// @Tags friends
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} models.Friend
// @Failure 500 {object} map[string]string
// @Router /friends [get]
func (h *UserHandler) GetFriends(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	friends, err := h.service.GetFriends(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, friends)
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	profile, err := h.service.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// RegisterUserRoutes sets up the routes for user handling with Gin
func RegisterUserRoutes(router *gin.Engine, userService *services.UserService) {
	handler := NewUserHandler(userService)

	userGroup := router.Group("/user")
	{
		userGroup.POST("/login", handler.Login)
		userGroup.POST("/register", handler.Register)
	}

	userProtectedGroup := userGroup
	userProtectedGroup.Use(middleware.JWTAuthMiddleware())
	{
		userProtectedGroup.GET("/profile", handler.GetProfile)
	}

	friendGroup := router.Group("/friends")
	friendGroup.Use(middleware.JWTAuthMiddleware())
	{
		friendGroup.POST("/:friend_id", handler.AddFriend)
		friendGroup.GET("", handler.GetFriends)
	}
}
