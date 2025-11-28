package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rizkylab/Go-SecureVibes/webui/backend/database"
)

// APIToken represents an API token
type APIToken struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Token     string     `json:"token,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	LastUsed  *time.Time `json:"last_used,omitempty"`
}

// CreateAPITokenRequest represents request to create API token
type CreateAPITokenRequest struct {
	Name      string `json:"name"`
	ExpiresIn int    `json:"expires_in"` // days, 0 = never expires
}

// ListAPITokens returns all API tokens
func ListAPITokens(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rows, err := db.Query(`
			SELECT id, name, created_at, expires_at, last_used
			FROM api_tokens
			ORDER BY created_at DESC
		`)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Failed to fetch API tokens",
			})
		}
		defer rows.Close()

		tokens := []APIToken{}
		for rows.Next() {
			var token APIToken
			var expiresAt, lastUsed *time.Time

			err := rows.Scan(&token.ID, &token.Name, &token.CreatedAt, &expiresAt, &lastUsed)
			if err != nil {
				continue
			}

			token.ExpiresAt = expiresAt
			token.LastUsed = lastUsed
			tokens = append(tokens, token)
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data":    tokens,
		})
	}
}

// CreateAPIToken creates a new API token
func CreateAPIToken(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req CreateAPITokenRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid request body",
			})
		}

		if req.Name == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "Token name is required",
			})
		}

		// Generate token
		tokenID := uuid.New().String()
		tokenValue := generateToken()

		var expiresAt *time.Time
		if req.ExpiresIn > 0 {
			expires := time.Now().AddDate(0, 0, req.ExpiresIn)
			expiresAt = &expires
		}

		// Insert token
		_, err := db.Exec(`
			INSERT INTO api_tokens (id, name, token, expires_at)
			VALUES (?, ?, ?, ?)
		`, tokenID, req.Name, tokenValue, expiresAt)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Failed to create API token",
			})
		}

		token := APIToken{
			ID:        tokenID,
			Name:      req.Name,
			Token:     tokenValue,
			CreatedAt: time.Now(),
			ExpiresAt: expiresAt,
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"success": true,
			"message": "API token created successfully. Save this token, it won't be shown again.",
			"data":    token,
		})
	}
}

// RevokeAPIToken revokes (deletes) an API token
func RevokeAPIToken(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenID := c.Params("id")

		result, err := db.Exec("DELETE FROM api_tokens WHERE id = ?", tokenID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Failed to revoke API token",
			})
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "API token not found",
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"message": "API token revoked successfully",
		})
	}
}

// generateToken generates a random API token
func generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return "sv_" + hex.EncodeToString(b)
}
