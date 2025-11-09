// Package handlers contains the delete handlers for the application.
package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/pageton/temp-mail/internal/db"
)

func DeleteInbox(c *fiber.Ctx) error {
	inboxID := c.Params("inboxid")
	if inboxID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Missing inbox ID")
	}
	queries := c.Locals("queries").(*db.Queries)
	inbox, err := queries.GetInboxByID(c.Context(), inboxID)
	if err != nil || inbox.ID == "" {
		log.Println("Error getting inbox:", err)
		return c.Status(fiber.StatusNotFound).
			JSON(&fiber.Map{"error": "Inbox does not exist or has been deleted already"})
	}

	err = queries.DeleteByInboxID(c.Context(), inbox.ID)
	if err != nil {
		log.Println("Error deleting inbox:", err)
		return c.Status(fiber.StatusBadRequest).
			JSON(&fiber.Map{"error": "Error deleting inbox"})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{"success": true})
}
