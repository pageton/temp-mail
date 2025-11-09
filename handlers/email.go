// Package handlers contains the get email handlers for the application.
package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/pageton/temp-mail/config"
	"github.com/pageton/temp-mail/internal/db"
)

type DatabaseEmail struct {
	ID          string    `json:"id"`
	Subject     *string   `json:"subject,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	ExpiresAt   time.Time `json:"expiresAt"`
	FromAddress *string   `json:"fromAddress,omitempty"`
	ToAddress   string    `json:"toAddress"`
}

type DatabaseEmails []DatabaseEmail

type EmailResponse struct {
	Success bool           `json:"success"`
	Data    DatabaseEmails `json:"data"`
}

func GetEmail(c *fiber.Ctx) error {
	email := c.Params("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Missing email")
	}
	cfg := c.Locals("config").(*config.Config)
	parts := strings.SplitN(email, "@", 2)
	if len(parts) != 2 {
		fmt.Println("Invalid email address")
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"error": "Invalid email address"})
	}
	domain := parts[1]
	if !slices.Contains(cfg.Domains.Aliases, domain) {
		log.Println("Email address does NOT belong to allowed domains:", domain)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email address does not belong to allowed domains",
		})
	}

	queries := c.Locals("queries").(*db.Queries)
	emails, err := queries.GetEmailsForAddress(
		c.Context(),
		sql.NullString{String: email, Valid: true},
	)
	if err != nil {
		log.Println("Error getting emails:", err)
		return c.Status(fiber.StatusInternalServerError).
			JSON(&fiber.Map{"error": "Error getting emails"})
	}

	if len(emails) == 0 {
		log.Println("No emails found for address:", email)
		return c.Status(fiber.StatusNotFound).
			JSON(&fiber.Map{"error": "No emails found for address"})
	}

	var result DatabaseEmails
	for _, e := range emails {
		de := DatabaseEmail{
			ID:          e.ID,
			Subject:     &e.Subject.String,
			CreatedAt:   time.UnixMilli(e.Createdat.Int64),
			ExpiresAt:   e.Expiresat.Time,
			FromAddress: &e.Fromaddress.String,
			ToAddress:   e.Toaddress,
		}
		result = append(result, de)
	}

	return c.Status(fiber.StatusOK).JSON(&EmailResponse{
		Success: true,
		Data:    result,
	})
}
