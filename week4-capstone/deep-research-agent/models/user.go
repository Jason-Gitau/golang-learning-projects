package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user account in the system
type User struct {
	ID        string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Email     string         `gorm:"unique;not null;type:varchar(255)" json:"email"`
	Password  string         `gorm:"not null;type:varchar(255)" json:"-"` // Never expose in JSON
	Name      string         `gorm:"not null;type:varchar(255)" json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Sessions []ResearchSession `gorm:"foreignKey:UserID" json:"-"`
	Documents []Document       `gorm:"foreignKey:UserID" json:"-"`
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}

// UserRegisterRequest represents the registration request payload
type UserRegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required,min=2,max=100"`
}

// UserLoginRequest represents the login request payload
type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserUpdateRequest represents the profile update request
type UserUpdateRequest struct {
	Name  string `json:"name" binding:"omitempty,min=2,max=100"`
	Email string `json:"email" binding:"omitempty,email"`
}

// ChangePasswordRequest represents the change password request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

// UserResponse represents the user data returned to clients (without sensitive data)
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// AuthResponse represents the authentication response with token
type AuthResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

// ToResponse converts a User to UserResponse (excludes password)
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
	}
}

// RateLimit tracks API rate limiting for users
type RateLimit struct {
	ID              string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID          string    `gorm:"not null;index;type:varchar(36)" json:"user_id"`
	RequestCount    int       `gorm:"not null;default:0" json:"request_count"`
	WindowStart     time.Time `gorm:"not null" json:"window_start"`
	ActiveJobs      int       `gorm:"not null;default:0" json:"active_jobs"`
	LastRequestAt   time.Time `json:"last_request_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// TableName specifies the table name for the RateLimit model
func (RateLimit) TableName() string {
	return "rate_limits"
}

// BeforeCreate hook to set ID using UUID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		// ID will be set by the handler using UUID
		return nil
	}
	return nil
}

// BeforeCreate hook for RateLimit
func (r *RateLimit) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		// ID will be set by the handler using UUID
		return nil
	}
	return nil
}
