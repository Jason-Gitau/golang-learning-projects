package utils

import (
	"errors"
	"fmt"

	"deep-research-agent/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	// ErrUserNotAuthenticated is returned when user is not in context
	ErrUserNotAuthenticated = errors.New("user not authenticated")
	// ErrUserNotFound is returned when user is not found in database
	ErrUserNotFound = errors.New("user not found")
)

// GetUserIDFromContext retrieves the user ID from the Gin context
func GetUserIDFromContext(c *gin.Context) (string, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", ErrUserNotAuthenticated
	}

	id, ok := userID.(string)
	if !ok || id == "" {
		return "", ErrUserNotAuthenticated
	}

	return id, nil
}

// GetUserFromContext retrieves the full user object from database using context
func GetUserFromContext(c *gin.Context, db *gorm.DB) (*models.User, error) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return nil, err
	}

	var user models.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	return &user, nil
}

// GetUserEmail retrieves the user email from context
func GetUserEmail(c *gin.Context) (string, error) {
	email, exists := c.Get("user_email")
	if !exists {
		return "", ErrUserNotAuthenticated
	}

	emailStr, ok := email.(string)
	if !ok || emailStr == "" {
		return "", ErrUserNotAuthenticated
	}

	return emailStr, nil
}

// GetUserName retrieves the user name from context
func GetUserName(c *gin.Context) (string, error) {
	name, exists := c.Get("user_name")
	if !exists {
		return "", ErrUserNotAuthenticated
	}

	nameStr, ok := name.(string)
	if !ok {
		return "", ErrUserNotAuthenticated
	}

	return nameStr, nil
}

// RequireAuth ensures user is authenticated and returns user ID
func RequireAuth(c *gin.Context) (string, error) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return "", err
	}
	return userID, nil
}

// IsAuthenticated checks if a user is authenticated
func IsAuthenticated(c *gin.Context) bool {
	_, err := GetUserIDFromContext(c)
	return err == nil
}

// GetUserInfo returns user information from context
type UserInfo struct {
	ID    string
	Email string
	Name  string
}

// GetUserInfo retrieves all user info from context
func GetUserInfo(c *gin.Context) (*UserInfo, error) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return nil, err
	}

	email, _ := GetUserEmail(c)
	name, _ := GetUserName(c)

	return &UserInfo{
		ID:    userID,
		Email: email,
		Name:  name,
	}, nil
}

// SetUserContext sets user information in the context
func SetUserContext(c *gin.Context, userID, email, name string) {
	c.Set("user_id", userID)
	c.Set("user_email", email)
	c.Set("user_name", name)
}
