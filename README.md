# Temp-Mail

A lightweight temporary email service built with Go, Fiber, and SQLite.

<p align="center">
    <img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License">
    <img src="https://img.shields.io/badge/go-1.20+-00ADD8.svg" alt="Go Version">
    <img src="https://img.shields.io/badge/fiber-v2.52.9+-green.svg" alt="Fiber Version">
    <img src="https://img.shields.io/badge/postfix-3.8.6-blue.svg" alt="Postfix Version">
</p>

A self-hosted temporary email service that provides disposable email addresses with automatic cleanup. Built for privacy and performance.

## Features

- **Fast & Lightweight**: Built with Go and Fiber v2 for high performance
- **Email Processing**: Full MIME email parsing with content extraction
- **SQLite Database**: Embedded database with automatic cleanup
- **Secure**: Webhook authentication and CORS protection
- **REST API**: Clean API endpoints for integration
- **Auto-cleanup**: Automatic deletion of expired emails
- **Easy Setup**: Minimal configuration required

## Installation

### Prerequisites

- Go 1.20+ or higher
- SQLite (included with Go)
- Postfix and procmail (for email processing)

### Quick Start

```bash
# Clone the repository
git clone https://github.com/pageton/temp-mail.git
cd temp-mail

# Install Go dependencies
go mod tidy

# Run the application
go run cmd/main.go
```
## Configuration

The application is configured via `config.toml`:

```toml
[app]
name = "temp-mail"
version = "0.1.0"

[server]
host = "localhost"
port = 3000
secret = 31  # Change this to a secure secret
prefork = false

[domains]
aliases = ["example.com", "example2.org"]  # Your domains

[database]
path = "mail.db"
```

### Email Setup (Production)

**For full email processing functionality:**

1. **DNS Configuration (Cloudflare)**:
   1. Log in to your Cloudflare account
   2. Select your domain (pageton.org)

      **Note**: In all the following steps, replace `pageton.org` with your domain name.

   3. Go to the DNS settings
   4. Add the following records:
      - MX record:
         - Name: `@`
         - Value: `mail.pageton.org`
         - Priority: 10
      - A record:
         - Name: `mail`
         - Value: `[Your VPS IP address]`
      - TXT record (for SPF):
         - Name: `@`
         - Value: `v=spf1 mx ~all`

2. **Install Postfix and procmail**:
   ```bash
   sudo apt update && sudo apt upgrade -y
   sudo apt install postfix procmail
   ```

3. **Set up domains**:
   - Add your domains to the `aliases` array in `config.toml`
   - Ensure Postfix is configured to accept these domains

The application will automatically configure Postfix on startup if the required system dependencies are available.

## API Documentation

### Endpoints

#### Get Available Domains
```http
GET /api/domains
```
Returns the list of configured domain aliases.

#### Get Emails for Address
```http
GET /api/email/:email
```
Retrieves all emails for a specific email address.

#### Get Individual Inbox
```http
GET /api/inbox/:inboxid
```
Fetches the full content of a specific inbox.

#### Delete Inbox
```http
GET /api/delete/:inboxid
```
Deletes an inbox and all associated emails.

#### Webhook Endpoint
```http
POST /webhook
```
Receives incoming emails from Postfix. Requires authentication via the secret header.

### Example Usage

```bash
# Get available domains
curl http://localhost:3000/api/domains

# Get emails for an address
curl http://localhost:3000/api/email/test@example.com

# Get specific inbox content
curl http://localhost:3000/api/inbox/inbox-id-here
```

## Development

### Project Structure

```
temp-mail/
├── cmd/main.go              # Application entry point
├── config/                  # Configuration handling
├── handlers/                # HTTP request handlers
├── internal/
│   ├── db/                  # Database layer (SQLC-generated)
│   ├── postfix/             # Postfix integration
│   ├── sqlc/                # SQL schemas and queries
│   └── utils/               # Utility functions
├── middlewares/             # Fiber middleware
└── config.toml              # Configuration file
```

### Database Operations

```bash
# Regenerate SQLC code after modifying queries.sql
sqlc generate

# The SQLite database is automatically created at "mail.db"
# Schema is applied automatically from embedded SQL
```

### Database Schema

The application uses three interconnected tables:

- **Email**: Stores email metadata with automatic expiration
- **Inbox**: Stores email content (text/HTML) linked to emails
- **EmailAddress**: Stores sender/recipient addresses

## Security

- Webhook authentication using configurable secret
- CORS protection for cross-origin requests
- Automatic email expiration (default: 3 days)
- SQLite foreign key constraints for data integrity

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

If you encounter any issues or have questions, please open an issue on the GitHub repository.
