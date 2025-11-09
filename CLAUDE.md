# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Building and Running
```bash
# Run the application
go run cmd/main.go

# Build the binary
go build -o temp-mail cmd/main.go
```

### Database Operations
```bash
# The SQLite database will be automatically created at "mail.db" on first run
# Schema is applied automatically from embedded SQL in internal/sqlc/schema.sql

# Regenerate SQLC code after modifying queries.sql
sqlc generate
```

### System Dependencies
```bash
# For production email processing (requires sudo)
sudo apt install postfix procmail

# Note: Without these system dependencies, the app can run but won't receive actual emails
```

### Go Dependencies
```bash
# Install dependencies
go mod tidy

# Download dependencies
go mod download
```

## Architecture Overview

This is a Go-based temporary email service built with the Fiber web framework and SQLite database.

### Core Architecture

**Web Framework**: Fiber v2 for HTTP server and routing
**Database**: SQLite with SQLC for type-safe database operations
**Email Processing**: enmime library for MIME email parsing
**Configuration**: TOML-based configuration system

### Key Components

1. **Main Application** (`cmd/main.go`):
   - Initializes Fiber app
   - Sets up Postfix email processing (fails gracefully if not available)
   - Loads configuration from `config.toml`
   - Sets up SQLite database with embedded schema
   - Configures database pragmas (WAL mode, foreign keys)
   - Sets up middleware and routes
   - Starts automatic cleanup task for expired emails

2. **Database Layer** (`internal/db/`):
   - SQLC-generated type-safe database operations
   - Three main tables: Email, Inbox, EmailAddress
   - Automatic schema creation from embedded SQL

3. **Handlers** (`handlers/`):
   - `domains.go`: Returns configured domain aliases
   - `webhook.go`: Processes incoming emails via webhook
   - `email.go`: Retrieves emails for a specific email address
   - `inbox.go`: Fetches individual inbox content by ID
   - `delete.go`: Deletes inboxes and associated emails

4. **Postfix Integration** (`internal/postfix/`):
   - Automatic Postfix configuration for email routing
   - Template-based configuration file generation
   - System service management for production deployment

5. **Configuration** (`config.toml`):
   - Server settings (host, port, webhook secret)
   - Domain aliases for temporary emails
   - Database path and logging configuration

6. **Utilities** (`internal/utils/`):
   - Email parsing and validation helpers
   - Automatic cleanup of expired emails

### Database Schema

The application uses three interconnected tables:
- **Email**: Stores email metadata with automatic expiration (3 days default)
- **Inbox**: Stores email content (text/HTML) with foreign key to Email
- **EmailAddress**: Stores sender/recipient addresses linked to emails

### Request Flow

1. Middleware injects config and database queries into request context
2. `GET /api/domains`: Returns available domain aliases
3. `POST /webhook`: Receives emails, validates webhook secret, parses MIME content, stores in database
4. `GET /api/email/:email`: Retrieves all emails for a specific email address
5. `GET /api/inbox/:inboxid`: Fetches individual inbox content by ID
6. `GET /api/delete/:inboxid`: Deletes inbox and associated emails

### Configuration

The application uses `config.toml` for all configuration. The file contains:
- App metadata and debug settings
- Server host, port, and webhook secret
- Logging level and output file
- Domain aliases for temporary email addresses
- Database file path

### Email Processing

The webhook handler:
1. Validates the webhook secret from headers
2. Parses incoming MIME emails using enmime
3. Extracts sender, recipients, subject, and body content
4. Stores email data with automatic expiration timestamps
5. Creates inbox entries for each recipient

### Important Notes

- Database file (`mail.db`) is ignored by git
- Webhook authentication relies on secret header validation
- Foreign key constraints are enforced for data integrity
- CORS middleware is configured for cross-origin requests
- The application handles graceful shutdown on SIGINT/SIGTERM
- **Postfix Setup**: The application requires Postfix and procmail for full email processing functionality. These are system dependencies that need sudo privileges to install. Without them, the app can run but won't receive actual emails.