package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gofiber/fiber/v2"
	_ "github.com/mattn/go-sqlite3"

	"github.com/pageton/temp-mail/config"
	"github.com/pageton/temp-mail/handlers"
	"github.com/pageton/temp-mail/internal/db"
	"github.com/pageton/temp-mail/internal/sqlc"
)

func main() {
	app := fiber.New()

	ctx := context.Background()
	cfg, err := config.LoadConfig("config.toml")
	if err != nil {
		log.Fatal(err)
	}

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

	queries := db.New(database)

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("config", cfg)
		c.Locals("queries", queries)
		return c.Next()
	})

	app.Get("/api/domains", handlers.GetDomains)

	app.Post("/webhook", handlers.Webhook)

	log.Fatal(app.Listen(cfg.Server.Host + ":" + strconv.Itoa(cfg.Server.Port)))
}
