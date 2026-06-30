package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kjanat/chatlogger-api-go/internal/domain"
	"github.com/kjanat/chatlogger-api-go/internal/hash"
	pkgerrors "github.com/pkg/errors"
)

// UserService implements the domain.UserService interface.
type UserService struct {
	userRepo  domain.UserRepository
	jwtSecret string
}

// NewUserService creates a new user service.
func NewUserService(userRepo domain.UserRepository, jwtSecret string) domain.UserService {
	return &UserService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// Authenticate authenticates a user and returns the user and JWT token.
func (s *UserService) Authenticate(email, password string) (*domain.User, string, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, "", fmt.Errorf("error finding user: %w", err)
	}

	if user == nil {
		return nil, "", errors.New("invalid email or password")
	}

	// Check password using our centralized hash package
	if err := hash.VerifyPassword(user.PasswordHash, password); err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := generateJWT(user, s.jwtSecret)
	if err != nil {
		return nil, "", fmt.Errorf("error generating token: %w", err)
	}

	// Update last login timestamp
	now := time.Now()
	user.LastLoginAt = &now
	if err := s.userRepo.Update(user); err != nil {
		// Non-critical error, log it but continue
		fmt.Printf("Error updating last login time: %v\n", err)
	}

	return user, token, nil
}

// Register registers a new user.
func (s *UserService) Register(user *domain.User, password string) error {
	// Check if user with this email already exists
	existingUser, err := s.userRepo.FindByEmail(user.Email)
	if err != nil {
		return fmt.Errorf("error checking existing user: %w", err)
	}

	if existingUser != nil {
		return errors.New("user with this email already exists")
	}

	// Hash the password using our centralized hash package
	hashedPassword, err := hash.GeneratePasswordHash(password, 10)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	// Set the hashed password and create timestamps
	user.PasswordHash = hashedPassword
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Create the user
	if err := s.userRepo.Create(user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetByID gets a user by ID.
func (s *UserService) GetByID(id uint64) (*domain.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by ID: %w", err)
	}
	return user, nil
}

// GetByEmail gets a user by email.
func (s *UserService) GetByEmail(email string) (*domain.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}
	return user, nil
}

// GetByOrganizationID gets users by organization ID with pagination.
func (s *UserService) GetByOrganizationID(orgID uint64, limit, offset int) ([]domain.User, error) {
	users, err := s.userRepo.FindByOrganizationID(orgID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to find users by organization ID: %w", err)
	}
	return users, nil
}

// UpdateUser updates a user.
func (s *UserService) UpdateUser(user *domain.User) error {
	// Set updated time
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// ChangePassword changes a user's password.
func (s *UserService) ChangePassword(userID uint64, currentPassword, newPassword string) error {
	// Get the user
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return fmt.Errorf("error finding user: %w", err)
	}

	if user == nil {
		return errors.New("user not found")
	}

	// Verify current password using our centralized hash package
	if err := hash.VerifyPassword(user.PasswordHash, currentPassword); err != nil {
		return errors.New("current password is incorrect")
	}

	// Hash the new password using our centralized hash package
	hashedPassword, err := hash.GeneratePasswordHash(newPassword, 10)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	// Update the user's password
	user.PasswordHash = hashedPassword
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to update user password: %w", err)
	}
	return nil
}

// DeleteUser deletes a user.
func (s *UserService) DeleteUser(id uint64) error {
	if err := s.userRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// JWTClaims represents the claims in a JWT token.
type JWTClaims struct {
	UserID         uint64      `json:"uid"`
	Email          string      `json:"email"`
	OrganizationID uint64      `json:"org_id"`
	Role           domain.Role `json:"role"`
	jwt.RegisteredClaims
}

// generateJWT generates a JWT token for a user.
func generateJWT(user *domain.User, secret string) (string, error) {
	// Set expiration time
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create claims
	claims := &JWTClaims{
		UserID:         user.ID,
		Email:          user.Email,
		OrganizationID: user.OrganizationID,
		Role:           user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", pkgerrors.Wrap(err, "failed to sign JWT token")
	}

	return tokenString, nil
}
