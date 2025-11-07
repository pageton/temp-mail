// Package handlers contains the GetDomains handlers for the application.
package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/pageton/temp-mail/config"
)

type Response struct {
	Success bool     `json:"success"`
	Result  []string `json:"result"`
}

func GetDomains(c *fiber.Ctx) error {
	cfg := c.Locals("config").(*config.Config)

	res := Response{
		Success: true,
		Result:  cfg.Domains.Aliases,
	}

	return c.Status(fiber.StatusOK).JSON(&res)
}
