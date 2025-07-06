package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kjanat/chatlogger-api-go/internal/domain"
)

// AuthHandler handles authentication-related requests.
type AuthHandler struct {
	userService domain.UserService
}

// NewAuthHandler creates a new authentication handler.
func NewAuthHandler(userService domain.UserService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
	}
}

// LoginRequest represents the login request body.
type LoginRequest struct {
	Email    string `binding:"required,email" json:"email"`
	Password string `binding:"required"       json:"password"`
}

// Login handles user login.
//
//	@Summary		User Login
//	@Description	Authenticates a user with email and password, returns user info and sets an HTTP-only auth cookie.
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		LoginRequest			true	"Login Credentials"
//	@Success		200		{object}	map[string]interface{}	"message: Login successful, user: domain.User"
//	@Failure		400		{object}	map[string]string		"Invalid request data"
//	@Failure		401		{object}	map[string]string		"Invalid email or password"
//	@Router			/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})

		return
	}

	// Authenticate user
	user, token, err := h.userService.Authenticate(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})

		return
	}

	// Set JWT token as HTTP-only cookie
	c.SetCookie(
		"auth_token",
		token,
		int(24*time.Hour.Seconds()), // 24 hours expiry
		"/",
		"",
		false, // secure (should be true in production with HTTPS)
		true,  // HTTP-only
	)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": user})
}

// RegisterRequest represents the register request body.
type RegisterRequest struct {
	Email     string      `binding:"required,email" json:"email"`
	Password  string      `binding:"required,min=8" json:"password"`
	FirstName string      `                         json:"first_name"`
	LastName  string      `                         json:"last_name"`
	Role      domain.Role `                         json:"role"`
	OrgID     uint64      `                         json:"organization_id"` // Made optional, otherwise: `binding:"required"`
}

// Register handles user registration.
//
//	@Summary		User Registration
//	@Description	Registers a new user. In a real app, this might require admin privileges or an invitation system.
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		RegisterRequest			true	"User Registration Details"
//	@Success		201		{object}	map[string]interface{}	"message: User registered successfully, user_id: uint64"
//	@Failure		400		{object}	map[string]string		"Invalid request data or registration failed (e.g., email exists)"
//	@Router			/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})

		return
	}

	// Use default organization if not provided (for easier testing)
	if req.OrgID == 0 {
		req.OrgID = 2 // Default to the second organization (unassigned)
	}

	// Set default role if not provided
	if req.Role == "" {
		req.Role = domain.RoleUser // Default to regular user
	}

	// Create user object
	user := &domain.User{
		Email:          req.Email,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Role:           req.Role,
		OrganizationID: req.OrgID,
	}

	// Register user
	if err := h.userService.Register(user, req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user_id": user.ID})
}

// Logout handles user logout.
//
//	@Summary		User Logout
//	@Description	Logs out the current user by clearing the authentication cookie.
//	@Tags			Authentication
//	@Produce		json
//	@Success		200	{object}	map[string]string	"Logout successful"
//	@Router			/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// Clear the auth cookie
	c.SetCookie(
		"auth_token",
		"",
		-1, // Expire immediately
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
