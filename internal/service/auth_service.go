package service

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/ethandiaz/yad-back/internal/domain"
	"github.com/ethandiaz/yad-back/pkg/jwt"
)

type AuthService struct {
	userRepo      domain.UserRepository
	tokenRepo     domain.AuthTokenRepository
	jwtManager    *jwt.JWTManager
	tokenDuration time.Duration
}

func NewAuthService(userRepo domain.UserRepository, tokenRepo domain.AuthTokenRepository) *AuthService {
	return &AuthService{
		userRepo:      userRepo,
		tokenRepo:     tokenRepo,
		jwtManager:    jwt.NewJWTManager(),
		tokenDuration: 24 * time.Hour,
	}
}

type LoginCredentials struct {
	Email    string
	Password string
}

type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	UserID       string `json:"userId"`
	Email        string `json:"email"`
	ExpiresAt    int64  `json:"expiresAt"`
}

func (s *AuthService) Login(email, password string) (*AuthResponse, error) {
	// Find user by email
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Verify password
	if !s.VerifyPassword(user, password) {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Generate access token
	accessToken, err := s.jwtManager.GenerateToken(user.ID, user.Email, s.tokenDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token (longer duration)
	refreshToken, err := s.jwtManager.GenerateToken(user.ID, user.Email, 7*24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Store refresh token in database
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	authToken := domain.NewAuthToken(user.ID, refreshToken, "refresh", expiresAt)
	if err := s.tokenRepo.Create(authToken); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       user.ID,
		Email:        user.Email,
		ExpiresAt:    time.Now().Add(s.tokenDuration).Unix(),
	}, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (*AuthResponse, error) {
	// Verify refresh token
	claims, err := s.jwtManager.VerifyToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Check if token exists in database
	dbToken, err := s.tokenRepo.GetByToken(refreshToken)
	if err != nil || dbToken == nil {
		return nil, fmt.Errorf("refresh token not found or revoked")
	}

	// Generate new access token
	accessToken, err := s.jwtManager.GenerateToken(claims.UserID, claims.Email, s.tokenDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Get user info
	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       user.ID,
		Email:        user.Email,
		ExpiresAt:    time.Now().Add(s.tokenDuration).Unix(),
	}, nil
}

func (s *AuthService) Logout(userID string) error {
	return s.tokenRepo.DeleteByUserID(userID)
}

func (s *AuthService) VerifyToken(tokenString string) (*jwt.Claims, error) {
	return s.jwtManager.VerifyToken(tokenString)
}

func (s *AuthService) VerifyPassword(user *domain.User, password string) bool {
	hashedPassword := s.hashPassword(password)
	return user.Password == hashedPassword
}

func (s *AuthService) hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash)
}
