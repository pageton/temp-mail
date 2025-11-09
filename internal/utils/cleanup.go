// Package utils contains the utility functions for the application.
package utils

import (
	"context"
	"database/sql"
	"time"
)

func StartCleanupTicker(ctx context.Context, db *sql.DB, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				db.Exec("DELETE FROM Email WHERE expiresAt <= CURRENT_TIMESTAMP")
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}
