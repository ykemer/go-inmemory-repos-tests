package middleware

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
)

// Recover catches any panic in a downstream handler, logs the value and full
// stack trace, then writes an RFC 7807 application/problem+json 500 response
// so the server keeps running instead of crashing.
//
// An optional logf function can be passed to redirect log output (useful in
// tests). When omitted, log.Printf is used.
func Recover(logf ...func(format string, args ...any)) fiber.Handler {
	logger := log.Printf
	if len(logf) > 0 && logf[0] != nil {
		logger = logf[0]
	}

	return func(c *fiber.Ctx) (err error) {
		defer func() {
			r := recover()
			if r == nil {
				return
			}

			logger("PANIC recovered on %s %s: %v\n%s",
				c.Method(), c.Path(), r, debug.Stack())

			p := problem{
				Type:   fmt.Sprintf("https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status/%d", fiber.StatusInternalServerError),
				Title:  http.StatusText(fiber.StatusInternalServerError),
				Status: fiber.StatusInternalServerError,
				Detail: "An unexpected error occurred",
			}
			err = c.Status(fiber.StatusInternalServerError).JSON(p, "application/problem+json")
		}()

		return c.Next()
	}
}
