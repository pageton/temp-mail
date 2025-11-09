// Package handlers contains the get inbox handlers for the application.
package handlers

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/pageton/temp-mail/internal/db"
)

type InboxResponse struct {
	ID          string    `json:"id"`
	TextContent *string   `json:"textContent"`
	HTMLContent *string   `json:"htmlContent"`
	Subject     *string   `json:"subject"`
	ExpiresAt   time.Time `json:"expiresAt"`
	CreatedAt   time.Time `json:"createdAt"`
	FromAddress *string   `json:"fromAddress"`
	ToAddress   string    `json:"toAddress"`
}

func GetInbox(c *fiber.Ctx) error {
	inboxID := c.Params("inboxid")
	if inboxID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Missing inbox ID")
	}
	queries := c.Locals("queries").(*db.Queries)
	inbox, err := queries.GetInboxByID(c.Context(), inboxID)
	if err != nil || inbox.ID == "" {
		log.Println("Error getting inbox:", err)
		return c.Status(fiber.StatusNotFound).
			SendString("Inbox does not exist or has been deleted")
	}
	return c.Status(fiber.StatusOK).JSON(&InboxResponse{
		ID:          inbox.ID,
		TextContent: &inbox.Textcontent.String,
		HTMLContent: &inbox.Htmlcontent.String,
		Subject:     &inbox.Subject.String,
		CreatedAt:   time.UnixMilli(inbox.Createdat.Int64),
		ExpiresAt:   inbox.Expiresat.Time,
		FromAddress: &inbox.Fromaddress.String,
		ToAddress:   inbox.Toaddress,
	})
}
