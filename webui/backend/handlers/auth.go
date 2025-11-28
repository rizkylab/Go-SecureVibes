package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rizkylab/Go-SecureVibes/webui/backend/config"
	"github.com/rizkylab/Go-SecureVibes/webui/backend/database"
	"github.com/rizkylab/Go-SecureVibes/webui/backend/middleware"
	"golang.org/x/crypto/bcrypt"
)

// LoginRequest represents login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents login response
type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	User      UserInfo  `json:"user"`
}

// UserInfo represents user information
type UserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

// Login handles user authentication
func Login(db *database.DB, cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid request body",
			})
		}

		// Validate input
		if req.Username == "" || req.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "Username and password are required",
			})
		}

		// Query user from database
		var user struct {
			ID           string
			Username     string
			PasswordHash string
			Email        string
			Role         string
		}

		err := db.QueryRow(`
			SELECT id, username, password_hash, email, role 
			FROM users 
			WHERE username = ?
		`, req.Username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Email, &user.Role)

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid username or password",
			})
		}

		// Verify password
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid username or password",
			})
		}

		// Parse token expiry
		expiry, _ := time.ParseDuration(cfg.Auth.TokenExpiry)
		if expiry == 0 {
			expiry = 24 * time.Hour
		}

		expiresAt := time.Now().Add(expiry)

		// Create JWT token
		claims := &middleware.Claims{
			UserID:   user.ID,
			Username: user.Username,
			Role:     user.Role,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expiresAt),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(cfg.Auth.JWTSecret))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Failed to generate token",
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data": LoginResponse{
				Token:     tokenString,
				ExpiresAt: expiresAt,
				User: UserInfo{
					ID:       user.ID,
					Username: user.Username,
					Email:    user.Email,
					Role:     user.Role,
				},
			},
		})
	}
}

// CreateDefaultUser creates a default admin user if no users exist
func CreateDefaultUser(db *database.DB) error {
	// Check if any users exist
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil // Users already exist
	}

	// Generate random user ID
	id := generateID()

	// Hash default password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Insert default admin user
	_, err = db.Exec(`
		INSERT INTO users (id, username, password_hash, email, role)
		VALUES (?, ?, ?, ?, ?)
	`, id, "admin", string(hashedPassword), "admin@securevibes.local", "admin")

	return err
}

// generateID generates a random ID
func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
