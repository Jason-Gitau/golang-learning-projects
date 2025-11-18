package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Email    string `gorm:"uniqueIndex;not null" json:"email" binding:"required,email"`
	Password string `gorm:"not null" json:"-"`
	Name     string `gorm:"not null" json:"name" binding:"required"`

	// Tier for rate limiting (e.g., "free", "premium", "enterprise")
	Tier string `gorm:"default:free" json:"tier"`

	// Relationships
	Agents        []Agent        `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"agents,omitempty"`
	Conversations []Conversation `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"conversations,omitempty"`
	UsageLogs     []UsageLog     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

// UserRegisterRequest represents a user registration request
type UserRegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

// UserLoginRequest represents a user login request
type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserResponse represents a user response (without password)
type UserResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Tier      string    `json:"tier"`
	CreatedAt time.Time `json:"created_at"`
}

// AuthResponse represents an authentication response with token
type AuthResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

// ToResponse converts a User to a UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		Tier:      u.Tier,
		CreatedAt: u.CreatedAt,
	}
}
