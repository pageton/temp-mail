package main

import (
	"context"
	"database/sql"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	_ "github.com/mattn/go-sqlite3"

	"github.com/pageton/temp-mail/config"
	"github.com/pageton/temp-mail/handlers"
	"github.com/pageton/temp-mail/internal/sqlc"
)

func main() {
	app := fiber.New()
	cfg, err := config.LoadConfig("config.toml")
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	db, err := sql.Open("sqlite3", cfg.Database.Path)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	if _, err = db.ExecContext(ctx, sqlc.Schema); err != nil {
		log.Fatal(err)
	}

	// Enable foreign key constraints
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal(err)
	}
	// Enable WAL mode
	_, err = db.Exec("PRAGMA journal_mode = WAL;")
	if err != nil {
		log.Fatal(err)
	}

	// Enable synchronous mode
	_, err = db.Exec("PRAGMA synchronous = NORMAL;")
	if err != nil {
		log.Fatal(err)
	}

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("config", cfg)
		return c.Next()
	})

	app.Get("/api/domains", handlers.GetDomains)

	log.Fatal(app.Listen(cfg.Server.Host + ":" + strconv.Itoa(cfg.Server.Port)))
}
