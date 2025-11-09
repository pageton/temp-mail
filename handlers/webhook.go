// Package handlers contains the webhook handlers for the application.
package handlers

import (
	"database/sql"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jhillyerd/enmime/v2"
	"github.com/lucsky/cuid"

	"github.com/pageton/temp-mail/config"
	"github.com/pageton/temp-mail/internal/db"
	"github.com/pageton/temp-mail/internal/utils"
)

type WebhookResponse struct {
	Success bool  `json:"success"`
	Data    int64 `json:"data"`
}

func Webhook(c *fiber.Ctx) error {
	cfg := c.Locals("config").(*config.Config)
	secret := c.Get("Secret")
	s, err := strconv.Atoi(secret)
	if err != nil || s != cfg.Server.Secret {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}
	r := strings.NewReader(string(c.Body()))

	env, err := enmime.ReadEnvelope(r)
	if err != nil {
		log.Println("Error parsing email:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error parsing email")
	}
	from := env.GetHeader("From")
	if from == "" {
		log.Println("Error getting from address:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error getting from address")
	}
	to := env.GetHeader("To")
	if to == "" {
		log.Println("Error getting to address:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error getting to address")
	}
	subject := env.GetHeader("Subject")
	if subject == "" {
		log.Println("Error getting subject:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error getting subject")
	}
	textBody := env.Text
	if textBody == "" {
		log.Println("Error getting body:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error getting body")
	}
	htmlBody := env.HTML
	if htmlBody == "" {
		log.Println("Error getting html:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error getting html")
	}

	toAddresses := utils.ParseEmailAddresses(to)
	fromAddresses := utils.ParseEmailAddresses(from)
	queries := c.Locals("queries").(*db.Queries)
	emailID, err := queries.InsertEmail(
		c.Context(),
		db.InsertEmailParams{
			Subject:   sql.NullString{String: subject, Valid: true},
			Expiresat: sql.NullTime{Time: time.Now().Add(3 * 24 * time.Hour), Valid: true},
		},
	)
	if err != nil {
		log.Println("Error inserting email:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error inserting email")
	}
	recipientGroups := []struct {
		Type   string
		Values []string
	}{
		{Type: "from", Values: fromAddresses},
		{Type: "to", Values: toAddresses},
	}

	for _, group := range recipientGroups {
		for _, addr := range group.Values {
			err = queries.InsertEmailAddress(
				c.Context(),
				db.InsertEmailAddressParams{
					Emailid: sql.NullInt64{Int64: emailID, Valid: true},
					Type:    sql.NullString{String: group.Type, Valid: true},
					Address: sql.NullString{String: addr, Valid: true},
				},
			)
			if err != nil {
				log.Println("Error inserting email address:", err)
				return c.Status(fiber.StatusInternalServerError).
					SendString("Error inserting email address")
			}
		}
	}
	for _, toAddress := range toAddresses {
		err = queries.InsertInbox(
			c.Context(),
			db.InsertInboxParams{
				ID:          cuid.New(),
				Emailid:     sql.NullInt64{Int64: emailID, Valid: true},
				Address:     sql.NullString{String: toAddress, Valid: true},
				Textcontent: sql.NullString{String: textBody, Valid: true},
				Htmlcontent: sql.NullString{String: htmlBody, Valid: true},
			},
		)
		if err != nil {
			log.Println("Error inserting inbox:", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error inserting inbox")
		}
	}
	res := WebhookResponse{
		Success: true,
		Data:    emailID,
	}
	return c.Status(fiber.StatusOK).JSON(&res)
}
