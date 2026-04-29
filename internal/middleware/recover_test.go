package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// captureLog returns a logf func and a pointer to the accumulated log output.
func captureLog() (func(string, ...any), *strings.Builder) {
	var buf strings.Builder
	logf := func(format string, args ...any) {
		fmt.Fprintf(&buf, format, args...)
	}
	return logf, &buf
}

func newRecoverApp(logf func(string, ...any), handler fiber.Handler) *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler:          ProblemJSONErrorHandler,
	})
	app.Use(Recover(logf))
	app.Get("/panic-string", func(c *fiber.Ctx) error {
		panic("something went wrong")
	})
	app.Get("/panic-error", func(c *fiber.Ctx) error {
		panic(fmt.Errorf("db connection lost"))
	})
	app.Get("/panic-value", func(c *fiber.Ctx) error {
		panic(42)
	})
	if handler != nil {
		app.Get("/ok", handler)
	} else {
		app.Get("/ok", func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
		})
	}
	return app
}

// ── Happy path ────────────────────────────────────────────────────────────────

func TestRecover_NoPanic_PassesThrough(t *testing.T) {
	logf, buf := captureLog()
	app := newRecoverApp(logf, nil)

	resp, err := app.Test(httptest.NewRequest("GET", "/ok", nil))

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Empty(t, buf.String(), "nothing should be logged when there is no panic")
}

// ── Panic types ───────────────────────────────────────────────────────────────

func TestRecover_PanicString_Returns500(t *testing.T) {
	logf, buf := captureLog()
	app := newRecoverApp(logf, nil)

	resp, err := app.Test(httptest.NewRequest("GET", "/panic-string", nil))

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, "application/problem+json", resp.Header.Get("Content-Type"))

	body, _ := io.ReadAll(resp.Body)
	var p map[string]any
	assert.NoError(t, json.Unmarshal(body, &p))
	assert.Equal(t, float64(500), p["status"])
	assert.Equal(t, "Internal Server Error", p["title"])
	assert.Equal(t, "An unexpected error occurred", p["detail"])

	assert.Contains(t, buf.String(), "something went wrong")
}

func TestRecover_PanicError_Returns500(t *testing.T) {
	logf, buf := captureLog()
	app := newRecoverApp(logf, nil)

	resp, err := app.Test(httptest.NewRequest("GET", "/panic-error", nil))

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	assert.Contains(t, buf.String(), "db connection lost")
}

func TestRecover_PanicArbitraryValue_Returns500(t *testing.T) {
	logf, _ := captureLog()
	app := newRecoverApp(logf, nil)

	resp, err := app.Test(httptest.NewRequest("GET", "/panic-value", nil))

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

// ── Log content ───────────────────────────────────────────────────────────────

func TestRecover_LogContainsMethodAndPath(t *testing.T) {
	logf, buf := captureLog()
	app := newRecoverApp(logf, nil)

	app.Test(httptest.NewRequest("GET", "/panic-string", nil))

	log := buf.String()
	assert.Contains(t, log, "GET")
	assert.Contains(t, log, "/panic-string")
}

func TestRecover_LogContainsStackTrace(t *testing.T) {
	logf, buf := captureLog()
	app := newRecoverApp(logf, nil)

	app.Test(httptest.NewRequest("GET", "/panic-string", nil))

	// runtime/debug.Stack always includes "goroutine" in its output
	assert.Contains(t, buf.String(), "goroutine")
}

// ── Default logger (no logf arg) ──────────────────────────────────────────────

func TestRecover_DefaultLogger_DoesNotPanic(t *testing.T) {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler:          ProblemJSONErrorHandler,
	})
	app.Use(Recover()) // no logf — falls back to log.Printf
	app.Get("/boom", func(c *fiber.Ctx) error {
		panic("test default logger")
	})

	resp, err := app.Test(httptest.NewRequest("GET", "/boom", nil))

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

// ── Server keeps running after panic ─────────────────────────────────────────

func TestRecover_ServerKeepsRunning(t *testing.T) {
	logf, _ := captureLog()
	app := newRecoverApp(logf, nil)

	// First request panics
	resp1, _ := app.Test(httptest.NewRequest("GET", "/panic-string", nil))
	assert.Equal(t, fiber.StatusInternalServerError, resp1.StatusCode)

	// Second request to a healthy endpoint still works
	resp2, _ := app.Test(httptest.NewRequest("GET", "/ok", nil))
	assert.Equal(t, fiber.StatusOK, resp2.StatusCode)
}
