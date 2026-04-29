package middleware

import (
	"tests/internal/idempotency"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Idempotency returns a Fiber middleware that deduplicates POST and PUT requests.
//
// Clients include an "Idempotency-Key" request header (typically a UUID).
// On the first request the response is processed normally and cached in store.
// Subsequent requests with the same key, method, and path return the cached
// response immediately, with the "X-Idempotency-Replayed: true" header set.
//
// Requests without an Idempotency-Key header pass through untouched, so the
// header is optional — existing clients need no changes.
func Idempotency(store idempotency.IStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Method() != fiber.MethodPost && c.Method() != fiber.MethodPut {
			return c.Next()
		}

		key := c.Get("Idempotency-Key")
		if key == "" {
			return c.Next()
		}

		// Scope the key to method + path so the same UUID cannot collide
		// across different endpoints.
		scopedKey := c.Method() + ":" + c.Path() + ":" + key

		entry, err := store.Get(c.Context(), scopedKey)
		if err != nil {
			return err
		}

		if entry != nil {
			c.Set("Idempotency-Key", key)
			c.Set("X-Idempotency-Replayed", "true")
			c.Set(fiber.HeaderContentType, entry.ContentType)
			return c.Status(entry.StatusCode).Send(entry.Body)
		}

		if err := c.Next(); err != nil {
			return err
		}

		// Copy the body — FastHTTP may reuse the buffer after the handler returns.
		body := make([]byte, len(c.Response().Body()))
		copy(body, c.Response().Body())

		_ = store.Set(c.Context(), scopedKey, idempotency.Entry{
			StatusCode:  c.Response().StatusCode(),
			Body:        body,
			ContentType: string(c.Response().Header.ContentType()),
			CreatedAt:   time.Now(),
		})

		return nil
	}
}
