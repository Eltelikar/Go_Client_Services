package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-pg/pg/v10"
)

const (
	maxRetries = 5
	retryDelay = 2 * time.Second
)

func (ps *PostService) RetryFunc(ctx context.Context, op func(tx *pg.Tx) error) error {
	var err error

	for i := 0; i < maxRetries; i++ {
		err = ps.db.RunInTransaction(ctx, op)
		if err == nil {
			return nil
		}

		if !mayRetry(err) {
			return err
		}

		time.Sleep(time.Duration(i+1) * retryDelay)
	}

	return fmt.Errorf("operation failed after %d retries: %w", maxRetries, err)
}

func mayRetry(err error) bool {
	if err == nil {
		return true
	}

	errMsg := err.Error()
	if strings.Contains(errMsg, "timeout") ||
		strings.Contains(errMsg, "deadlock detected") ||
		strings.Contains(errMsg, "canceling statement due to conflict") ||
		strings.Contains(errMsg, "could not serialize access") {
		return false
	}
	return true
}
