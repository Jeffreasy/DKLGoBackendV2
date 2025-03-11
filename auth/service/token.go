package service

import (
	"dklautomationgo/models"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidJWT = errors.New("ongeldige JWT token")
)

// TokenService handelt JWT token generatie en validatie
type TokenService struct {
	secretKey []byte
}

// NewTokenService maakt een nieuwe TokenService
func NewTokenService() *TokenService {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		log.Println("[TokenService] WARNING: JWT_SECRET_KEY is not set, using default (insecure) key")
		secretKey = "default-insecure-jwt-secret-key-change-in-production"
	}

	return &TokenService{
		secretKey: []byte(secretKey),
	}
}

// Claims representeert de JWT claims
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateAccessToken genereert een JWT access token voor een gebruiker
func (s *TokenService) GenerateAccessToken(user *models.User) (string, error) {
	expirationTime := time.Now().Add(getAccessTokenExpiry())

	claims := &Claims{
		UserID: user.ID.String(),
		Email:  user.Email,
		Role:   string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "dklautomationgo",
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		log.Printf("[TokenService] Error signing token: %v", err)
		return "", err
	}

	return tokenString, nil
}

// ValidateToken valideert een JWT token en geeft de claims terug
func (s *TokenService) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Valideer signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("onverwachte signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		log.Printf("[TokenService] Error parsing token: %v", err)
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidJWT
	}

	return claims, nil
}

// GetUserIDFromToken haalt de gebruiker ID uit een gevalideerde token
func (s *TokenService) GetUserIDFromToken(tokenString string) (uuid.UUID, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return uuid.Nil, err
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		log.Printf("[TokenService] Error parsing user ID from token: %v", err)
		return uuid.Nil, err
	}

	return userID, nil
}
