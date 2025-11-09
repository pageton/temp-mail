package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/mattn/go-sqlite3"

	"github.com/pageton/temp-mail/config"
	"github.com/pageton/temp-mail/handlers"
	"github.com/pageton/temp-mail/internal/db"
	"github.com/pageton/temp-mail/internal/sqlc"
	"github.com/pageton/temp-mail/internal/utils"
	"github.com/pageton/temp-mail/middlewares"
)

func main() {
	cfg, err := config.LoadConfig("config.toml")
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New(fiber.Config{Prefork: cfg.Server.Prefork})

	ctx := context.Background()

	database, err := sql.Open("sqlite3", cfg.Database.Path)
	if err != nil {
		log.Fatal(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("Received signal %s, closing database", sig)
		database.Close()
		log.Println("Database connection closed")
		os.Exit(0)
	}()

	if _, err = database.ExecContext(ctx, sqlc.Schema); err != nil {
		log.Fatal(err)
	}

	_, err = database.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal(err)
	}

	_, err = database.Exec("PRAGMA journal_mode = WAL;")
	if err != nil {
		log.Fatal(err)
	}

	_, err = database.Exec("PRAGMA synchronous = NORMAL;")
	if err != nil {
		log.Fatal(err)
	}

	middlewares.Cors(app) // CORS middleware

	middlewares.RateLimiter(app) // Rate limiter middleware

	utils.StartCleanupTicker(ctx, database, time.Hour*2) // Cleanup ticker

	queries := db.New(database)

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("config", cfg)
		c.Locals("queries", queries)
		return c.Next()
	})

	app.Post("/webhook", handlers.Webhook)

	api := app.Group("/api")
	api.Get("/domains", handlers.GetDomains)
	api.Get("/delete/:inboxid", handlers.DeleteInbox)
	api.Get("/email/:email", handlers.GetEmail)
	api.Get("/inbox/:inboxid", handlers.GetInbox)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Fatal(app.Listen(addr))
}
