// Package middlewares contains the middlewares for the application.
package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Cors(app *fiber.App) {
	app.Use(cors.New())
}
