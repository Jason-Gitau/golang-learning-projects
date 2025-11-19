package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims represents the JWT claims structure
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

// JWTManager handles JWT token operations
type JWTManager struct {
	secret         []byte
	tokenExpiry    time.Duration
	refreshExpiry  time.Duration
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(secret string, tokenExpiry, refreshExpiry time.Duration) *JWTManager {
	return &JWTManager{
		secret:        []byte(secret),
		tokenExpiry:   tokenExpiry,
		refreshExpiry: refreshExpiry,
	}
}

// GenerateToken generates a new JWT access token for a user
func (j *JWTManager) GenerateToken(userID, email, name string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		Name:   name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "deep-research-agent",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func (j *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	// Check if token is expired
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

// RefreshToken generates a new token from a valid existing token
func (j *JWTManager) RefreshToken(tokenString string) (string, error) {
	// Validate the existing token
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		// Allow refresh if token is only expired, not invalid
		token, parseErr := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return j.secret, nil
		})

		if parseErr != nil {
			return "", fmt.Errorf("invalid token: %w", parseErr)
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			return "", errors.New("invalid token claims")
		}

		// Check if token was issued within refresh window
		if claims.IssuedAt == nil {
			return "", errors.New("token has no issue time")
		}

		issuedTime := claims.IssuedAt.Time
		if time.Since(issuedTime) > j.refreshExpiry {
			return "", errors.New("refresh window expired")
		}

		// Generate new token with same user info
		return j.GenerateToken(claims.UserID, claims.Email, claims.Name)
	}

	// Generate new token
	return j.GenerateToken(claims.UserID, claims.Email, claims.Name)
}

// ExtractToken extracts the token from Authorization header
func ExtractToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	// Check for "Bearer " prefix
	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		return "", errors.New("invalid authorization header format. Expected: Bearer <token>")
	}

	return authHeader[len(bearerPrefix):], nil
}
