// Package middleware provides HTTP middleware for the ChatLogger API.
// This file implements authentication middleware including JWT-based auth for
// the dashboard UI and API key-based auth for external integrations.
package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kjanat/chatlogger-api-go/internal/domain"
	"github.com/kjanat/chatlogger-api-go/internal/service"
)

// Context keys for values stored in Gin context.
const (
	// OrganizationIDKey is the key used to store organization ID in the context.
	OrganizationIDKey = "orgID"
	// UserIDKey is the key used to store user ID in the context.
	UserIDKey = "userID"
	// RoleKey is the key used to store user role in the context.
	RoleKey = "role"
	// RequestedOrgIDKey is the key used to store requested organization ID in the context.
	RequestedOrgIDKey = "requestedOrgID"
)

// JWTAuth middleware for user authentication using JWT.
func JWTAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the JWT token from the cookie
		tokenString, err := c.Cookie("auth_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			c.Abort()

			return
		}

		// Parse the JWT token
		token, err := jwt.ParseWithClaims(
			tokenString,
			&service.JWTClaims{},
			func(token *jwt.Token) (interface{}, error) {
				// Validate the signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}

				return []byte(jwtSecret), nil
			},
		)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()

			return
		}

		// Get claims from token
		claims, ok := token.Claims.(*service.JWTClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()

			return
		}

		// Set user details in context
		c.Set(UserIDKey, claims.UserID)
		c.Set(OrganizationIDKey, claims.OrganizationID)
		c.Set(RoleKey, claims.Role)

		c.Next()
	}
}

// APIKeyAuth middleware for authentication using API key.
func APIKeyAuth(apiKeyService domain.APIKeyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the API key from the header
		apiKey := c.GetHeader("x-organization-api-key")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API key required"})
			c.Abort()

			return
		}

		// Validate the API key
		key, err := apiKeyService.ValidateKey(apiKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to validate API key"})
			c.Abort()

			return
		}

		if key == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid API key"})
			c.Abort()

			return
		}

		// Set organization ID in context
		c.Set(OrganizationIDKey, key.OrganizationID)

		c.Next()
	}
}

// RoleRequired middleware to check if user has required role.
func RoleRequired(roles ...domain.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user's role from context
		role, exists := c.Get(RoleKey)
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "role information not available"})
			c.Abort()

			return
		}

		userRole := role.(domain.Role)

		// Check if the user has one of the required roles
		authorized := false

		for _, r := range roles {
			if userRole == r {
				authorized = true

				break
			}
		}

		// SuperAdmin role has access to everything
		if userRole == domain.RoleSuperAdmin {
			authorized = true
		}

		if !authorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()

			return
		}

		c.Next()
	}
}

// ValidateOrgAccess middleware to check if user has access to the requested organization.
func ValidateOrgAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user's organization ID from context
		userOrgID, exists := c.Get(OrganizationIDKey)
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "organization information not available"})
			c.Abort()

			return
		}

		// Get requested organization ID from URL
		// This assumes the URL has a parameter named "orgID"
		requestedOrgIDStr := c.Param("orgID")
		if requestedOrgIDStr == "me" {
			// The "me" shorthand refers to the user's own organization
			c.Next()

			return
		}

		var requestedOrgID uint64
		if _, err := fmt.Sscanf(requestedOrgIDStr, "%d", &requestedOrgID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization ID"})
			c.Abort()

			return
		}

		// Get user's role from context
		role, exists := c.Get(RoleKey)
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "role information not available"})
			c.Abort()

			return
		}

		userRole := role.(domain.Role)

		// SuperAdmin can access any organization
		if userRole == domain.RoleSuperAdmin {
			c.Next()

			return
		}

		// Other users can only access their own organization
		if userOrgID.(uint64) != requestedOrgID {
			c.JSON(
				http.StatusForbidden,
				gin.H{"error": "you do not have access to this organization"},
			)
			c.Abort()

			return
		}

		c.Next()
	}
}

// ValidateSlugAccess middleware to check if user has access to the requested organization by slug.
func ValidateSlugAccess(orgService domain.OrganizationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user's organization ID from context
		userOrgID, exists := c.Get(OrganizationIDKey)
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "organization information not available"})
			c.Abort()

			return
		}

		// Get requested organization slug from URL
		slug := c.Param("slug")

		// Lookup the organization by slug
		org, err := orgService.GetBySlug(slug)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to lookup organization"})
			c.Abort()

			return
		}

		if org == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "organization not found"})
			c.Abort()

			return
		}

		// Store the organization ID in context for later use
		c.Set(RequestedOrgIDKey, org.ID)

		// Get user's role from context if it exists (may not exist for API key auth)
		roleInterface, roleExists := c.Get(RoleKey)

		// If role exists, check permissions
		if roleExists {
			userRole := roleInterface.(domain.Role)

			// SuperAdmin can access any organization
			if userRole == domain.RoleSuperAdmin {
				c.Next()

				return
			}

			// Other users can only access their own organization
			if userOrgID.(uint64) != org.ID {
				c.JSON(
					http.StatusForbidden,
					gin.H{"error": "you do not have access to this organization"},
				)
				c.Abort()

				return
			}
		} else if userOrgID.(uint64) != org.ID {
			// For API key auth, we already verified the key belongs to the org
			// Just check that the key's org matches the requested org
			c.JSON(http.StatusForbidden, gin.H{"error": "this API key cannot access this organization"})
			c.Abort()

			return
		}

		c.Next()
	}
}
