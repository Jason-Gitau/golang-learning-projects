package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens and adds user info to context
func AuthMiddleware(jwtManager *JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Extract token from header
		tokenString, err := ExtractToken(authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		// Validate the token
		claims, err := jwtManager.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
				"detail": err.Error(),
			})
			c.Abort()
			return
		}

		// Add user info to context for use in handlers
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_name", claims.Name)
		c.Set("claims", claims)

		c.Next()
	}
}

// OptionalAuthMiddleware validates JWT tokens if present, but allows requests without auth
func OptionalAuthMiddleware(jwtManager *JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// No auth provided, continue without user context
			c.Next()
			return
		}

		// Extract token from header
		tokenString, err := ExtractToken(authHeader)
		if err != nil {
			// Invalid format, but since auth is optional, continue
			c.Next()
			return
		}

		// Validate the token
		claims, err := jwtManager.ValidateToken(tokenString)
		if err != nil {
			// Invalid token, but since auth is optional, continue
			c.Next()
			return
		}

		// Add user info to context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_name", claims.Name)
		c.Set("claims", claims)

		c.Next()
	}
}

// WebSocketAuthMiddleware validates JWT tokens from query parameters for WebSocket connections
func WebSocketAuthMiddleware(jwtManager *JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to get token from query parameter first (for WebSocket connections)
		tokenString := c.Query("token")

		// If not in query, check Authorization header
		if tokenString == "" {
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Authentication required. Provide token in query parameter or Authorization header",
				})
				c.Abort()
				return
			}

			var err error
			tokenString, err = ExtractToken(authHeader)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": err.Error(),
				})
				c.Abort()
				return
			}
		}

		// Validate the token
		claims, err := jwtManager.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
				"detail": err.Error(),
			})
			c.Abort()
			return
		}

		// Add user info to context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_name", claims.Name)
		c.Set("claims", claims)

		c.Next()
	}
}

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// GetUserIDFromContext retrieves the user ID from the Gin context
func GetUserIDFromContext(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}
	id, ok := userID.(string)
	return id, ok
}

// RequireAuth ensures a user is authenticated, returns error if not
func RequireAuth(c *gin.Context) (string, error) {
	userID, exists := GetUserIDFromContext(c)
	if !exists || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		c.Abort()
		return "", ErrUnauthorized
	}
	return userID, nil
}

// Custom errors
var (
	ErrUnauthorized = &AuthError{Message: "unauthorized"}
)

// AuthError represents an authentication error
type AuthError struct {
	Message string
}

func (e *AuthError) Error() string {
	return e.Message
}

// GetAuthHeader extracts the Authorization header value
func GetAuthHeader(c *gin.Context) string {
	return strings.TrimSpace(c.GetHeader("Authorization"))
}
