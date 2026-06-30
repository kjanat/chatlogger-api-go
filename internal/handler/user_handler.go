package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kjanat/chatlogger-api-go/internal/domain"
)

// UserHandler handles user-related requests.
type UserHandler struct {
	userService domain.UserService
}

// NewUserHandler creates a new user handler.
func NewUserHandler(userService domain.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetMe handles the request to get the current user's info.
//
//	@Summary		Get Current User
//	@Description	Retrieves the profile information for the currently authenticated user.
//	@Tags			Users
//	@Produce		json
//	@Success		200	{object}	domain.User			"Current user's profile"
//	@Failure		401	{object}	map[string]string	"Unauthorized (JWT invalid/missing or User ID not found)"
//	@Failure		404	{object}	map[string]string	"User not found"
//	@Failure		500	{object}	map[string]string	"Failed to get user"
//	@Security		BearerAuth
//	@Router			/v1/users/me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	// Get user ID from context (set by JWTAuth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})

		return
	}

	// Get user from service
	user, err := h.userService.GetByID(userID.(uint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})

		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})

		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateMeRequest represents the update user request body.
type UpdateMeRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// UpdateMe handles the request to update the current user's info.
//
//	@Summary		Update Current User
//	@Description	Updates the first name and last name for the currently authenticated user.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		UpdateMeRequest		true	"User fields to update"
//	@Success		200		{object}	domain.User			"Updated user profile"
//	@Failure		400		{object}	map[string]string	"Invalid request data"
//	@Failure		401		{object}	map[string]string	"Unauthorized (JWT invalid/missing or User ID not found)"
//	@Failure		404		{object}	map[string]string	"User not found"
//	@Failure		500		{object}	map[string]string	"Failed to get or update user"
//	@Security		BearerAuth
//	@Router			/v1/users/me [patch]
func (h *UserHandler) UpdateMe(c *gin.Context) {
	var req UpdateMeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})

		return
	}

	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})

		return
	}

	// Get the current user
	user, err := h.userService.GetByID(userID.(uint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})

		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})

		return
	}

	// Update the user fields
	user.FirstName = req.FirstName
	user.LastName = req.LastName

	// Save the updated user
	if err := h.userService.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})

		return
	}

	c.JSON(http.StatusOK, user)
}

// ChangePasswordRequest represents the change password request body.
type ChangePasswordRequest struct {
	CurrentPassword string `binding:"required"       json:"current_password"`
	NewPassword     string `binding:"required,min=8" json:"new_password"`
}

// ChangePassword handles the request to change the current user's password.
//
//	@Summary		Change Password
//	@Description	Allows the currently authenticated user to change their password.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		ChangePasswordRequest	true	"Current and new password"
//	@Success		200		{object}	map[string]string		"Password changed successfully"
//	@Failure		400		{object}	map[string]string		"Invalid request data or password change failed (e.g., wrong current password)"
//	@Failure		401		{object}	map[string]string		"Unauthorized (JWT invalid/missing or User ID not found)"
//	@Security		BearerAuth
//	@Router			/v1/users/me/password [post]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})

		return
	}

	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})

		return
	}

	// Change the password
	err := h.userService.ChangePassword(userID.(uint64), req.CurrentPassword, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

// ListOrgUsers handles the request to list all users in the current organization.
//
//	@Summary		List Organization Users (Admin)
//	@Description	Retrieves a paginated list of all users belonging to the authenticated user's organization. Requires admin role.
//	@Tags			Users (Admin)
//	@Produce		json
//	@Param			limit	query		int					false	"Number of users per page"	default(20)
//	@Param			offset	query		int					false	"Offset for pagination"		default(0)
//	@Success		200		{array}		domain.User			"List of users in the organization"
//	@Failure		401		{object}	map[string]string	"Unauthorized (JWT invalid/missing or Org ID not found)"
//	@Failure		403		{object}	map[string]string	"Forbidden (User does not have admin role)"	//	Assuming	RoleRequired	middleware	handles	this
//	@Failure		500		{object}	map[string]string	"Failed to get users"
//	@Security		BearerAuth
//	@Router			/v1/orgs/me/users [get] // Assuming this route exists and is admin-protected
func (h *UserHandler) ListOrgUsers(c *gin.Context) {
	// Get organization ID from context
	orgID, exists := c.Get("orgID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization ID not found in context"})

		return
	}

	// Parse pagination parameters
	limit := 20 // Default limit
	offset := 0 // Default offset

	// Get users
	users, err := h.userService.GetByOrganizationID(orgID.(uint64), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})

		return
	}

	c.JSON(http.StatusOK, users)
}
