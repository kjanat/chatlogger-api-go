package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kjanat/chatlogger-api-go/internal/domain"
	"github.com/kjanat/chatlogger-api-go/internal/service"
	"github.com/kjanat/chatlogger-api-go/test/fixtures"
	"github.com/kjanat/chatlogger-api-go/test/mocks"
	"github.com/kjanat/chatlogger-api-go/test/testutils"
	"github.com/stretchr/testify/assert"
)

func createTestJWT(userID, orgID uint64, role domain.Role, secret string) (string, error) {
	claims := &service.JWTClaims{
		UserID:         userID,
		Email:          "test@example.com",
		OrganizationID: orgID,
		Role:           role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func TestJWTAuth_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := testutils.SetupTestRouter()
	
	jwtSecret := "test-secret"
	tokenString, err := createTestJWT(1, 100, domain.RoleUser, jwtSecret)
	assert.NoError(t, err)
	
	router.Use(JWTAuth(jwtSecret))
	router.GET("/test", func(c *gin.Context) {
		userID, _ := c.Get(UserIDKey)
		orgID, _ := c.Get(OrganizationIDKey)
		role, _ := c.Get(RoleKey)
		
		c.JSON(http.StatusOK, gin.H{
			"user_id": userID,
			"org_id":  orgID,
			"role":    role,
		})
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  "auth_token",
		Value: tokenString,
	})
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"user_id":1`)
	assert.Contains(t, w.Body.String(), `"org_id":100`)
	assert.Contains(t, w.Body.String(), `"role":"user"`)
}

func TestJWTAuth_NoCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := testutils.SetupTestRouter()
	
	router.Use(JWTAuth("test-secret"))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "authentication required")
}

func TestJWTAuth_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := testutils.SetupTestRouter()
	
	router.Use(JWTAuth("test-secret"))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  "auth_token",
		Value: "invalid.jwt.token",
	})
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid token")
}

func TestJWTAuth_ExpiredToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := testutils.SetupTestRouter()
	
	jwtSecret := "test-secret"
	
	// Create expired token
	claims := &service.JWTClaims{
		UserID:         1,
		Email:          "test@example.com",
		OrganizationID: 100,
		Role:           domain.RoleUser,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Expired
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	assert.NoError(t, err)
	
	router.Use(JWTAuth(jwtSecret))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  "auth_token",
		Value: tokenString,
	})
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid token")
}

func TestAPIKeyAuth_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := testutils.SetupTestRouter()
	
	mockAPIKeyService := &mocks.MockAPIKeyService{}
	apiKey := &domain.APIKey{
		ID:             1,
		OrganizationID: 100,
		HashedKey:      "$2a$10$hashed.api.key",
		Label:          "Test Key",
	}
	
	mockAPIKeyService.On("ValidateKey", "test-api-key").Return(apiKey, nil)
	
	router.Use(APIKeyAuth(mockAPIKeyService))
	router.GET("/test", func(c *gin.Context) {
		orgID, _ := c.Get(OrganizationIDKey)
		c.JSON(http.StatusOK, gin.H{
			"org_id": orgID,
		})
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("x-organization-api-key", "test-api-key")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"org_id":100`)
	
	mockAPIKeyService.AssertExpectations(t)
}

func TestAPIKeyAuth_NoHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := testutils.SetupTestRouter()
	
	mockAPIKeyService := &mocks.MockAPIKeyService{}
	
	router.Use(APIKeyAuth(mockAPIKeyService))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "API key required")
}

func TestAPIKeyAuth_InvalidKey(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := testutils.SetupTestRouter()
	
	mockAPIKeyService := &mocks.MockAPIKeyService{}
	mockAPIKeyService.On("ValidateKey", "invalid-key").Return((*domain.APIKey)(nil), nil)
	
	router.Use(APIKeyAuth(mockAPIKeyService))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("x-organization-api-key", "invalid-key")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid API key")
	
	mockAPIKeyService.AssertExpectations(t)
}

func TestRoleRequired_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := testutils.SetupTestRouter()
	
	router.Use(func(c *gin.Context) {
		c.Set(RoleKey, domain.RoleAdmin)
		c.Next()
	})
	router.Use(RoleRequired(domain.RoleAdmin, domain.RoleUser))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoleRequired_SuperAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := testutils.SetupTestRouter()
	
	router.Use(func(c *gin.Context) {
		c.Set(RoleKey, domain.RoleSuperAdmin)
		c.Next()
	})
	router.Use(RoleRequired(domain.RoleAdmin)) // SuperAdmin should have access
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoleRequired_InsufficientPermissions(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := testutils.SetupTestRouter()
	
	router.Use(func(c *gin.Context) {
		c.Set(RoleKey, domain.RoleViewer)
		c.Next()
	})
	router.Use(RoleRequired(domain.RoleAdmin))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "insufficient permissions")
}

func TestRoleRequired_NoRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := testutils.SetupTestRouter()
	
	router.Use(RoleRequired(domain.RoleUser))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "role information not available")
}

func TestValidateOrgAccess_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := testutils.SetupTestRouter()
	
	router.Use(func(c *gin.Context) {
		c.Set(OrganizationIDKey, uint64(100))
		c.Set(RoleKey, domain.RoleUser)
		c.Next()
	})
	router.Use(ValidateOrgAccess())
	router.GET("/test/:orgID", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	req := httptest.NewRequest("GET", "/test/100", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestValidateOrgAccess_MeShorthand(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := testutils.SetupTestRouter()
	
	router.Use(func(c *gin.Context) {
		c.Set(OrganizationIDKey, uint64(100))
		c.Set(RoleKey, domain.RoleUser)
		c.Next()
	})
	router.Use(ValidateOrgAccess())
	router.GET("/test/:orgID", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	req := httptest.NewRequest("GET", "/test/me", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestValidateOrgAccess_SuperAdminAccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := testutils.SetupTestRouter()
	
	router.Use(func(c *gin.Context) {
		c.Set(OrganizationIDKey, uint64(100))
		c.Set(RoleKey, domain.RoleSuperAdmin)
		c.Next()
	})
	router.Use(ValidateOrgAccess())
	router.GET("/test/:orgID", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	req := httptest.NewRequest("GET", "/test/200", nil) // Different org
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestValidateOrgAccess_Forbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := testutils.SetupTestRouter()
	
	router.Use(func(c *gin.Context) {
		c.Set(OrganizationIDKey, uint64(100))
		c.Set(RoleKey, domain.RoleUser)
		c.Next()
	})
	router.Use(ValidateOrgAccess())
	router.GET("/test/:orgID", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	req := httptest.NewRequest("GET", "/test/200", nil) // Different org
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "you do not have access to this organization")
}

func TestValidateSlugAccess_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := testutils.SetupTestRouter()
	
	mockOrgService := &mocks.MockOrganizationService{}
	org := fixtures.CreateTestOrganization()
	org.ID = 100
	org.Slug = "test-org"
	
	mockOrgService.On("GetBySlug", "test-org").Return(org, nil)
	
	router.Use(func(c *gin.Context) {
		c.Set(OrganizationIDKey, uint64(100))
		c.Set(RoleKey, domain.RoleUser)
		c.Next()
	})
	router.Use(ValidateSlugAccess(mockOrgService))
	router.GET("/test/:slug", func(c *gin.Context) {
		requestedOrgID, _ := c.Get(RequestedOrgIDKey)
		c.JSON(http.StatusOK, gin.H{
			"message":           "success",
			"requested_org_id": requestedOrgID,
		})
	})
	
	req := httptest.NewRequest("GET", "/test/test-org", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"requested_org_id":100`)
	
	mockOrgService.AssertExpectations(t)
}

func TestValidateSlugAccess_OrgNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := testutils.SetupTestRouter()
	
	mockOrgService := &mocks.MockOrganizationService{}
	mockOrgService.On("GetBySlug", "nonexistent").Return((*domain.Organization)(nil), nil)
	
	router.Use(func(c *gin.Context) {
		c.Set(OrganizationIDKey, uint64(100))
		c.Set(RoleKey, domain.RoleUser)
		c.Next()
	})
	router.Use(ValidateSlugAccess(mockOrgService))
	router.GET("/test/:slug", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	req := httptest.NewRequest("GET", "/test/nonexistent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "organization not found")
	
	mockOrgService.AssertExpectations(t)
}

func TestValidateSlugAccess_APIKeyAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := testutils.SetupTestRouter()
	
	mockOrgService := &mocks.MockOrganizationService{}
	org := fixtures.CreateTestOrganization()
	org.ID = 100
	org.Slug = "test-org"
	
	mockOrgService.On("GetBySlug", "test-org").Return(org, nil)
	
	router.Use(func(c *gin.Context) {
		c.Set(OrganizationIDKey, uint64(100))
		// No role set (API key auth scenario)
		c.Next()
	})
	router.Use(ValidateSlugAccess(mockOrgService))
	router.GET("/test/:slug", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	req := httptest.NewRequest("GET", "/test/test-org", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	mockOrgService.AssertExpectations(t)
}