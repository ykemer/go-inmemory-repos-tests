package middleware

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"tests/internal/idempotency"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// newTestApp builds a minimal Fiber app wired with the idempotency middleware.
func newTestApp(store idempotency.IStore) *fiber.App {
	app := fiber.New(fiber.Config{
		// Suppress startup/shutdown logs in tests.
		DisableStartupMessage: true,
	})
	app.Use(Idempotency(store))

	app.Post("/resources", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"created": true})
	})
	app.Put("/resources/:id", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"updated": true})
	})
	app.Get("/resources", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"list": true})
	})
	app.Delete("/resources/:id", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNoContent).SendString("")
	})
	return app
}

func postWithKey(app *fiber.App, key string) (*httptest.ResponseRecorder, error) {
	req := httptest.NewRequest("POST", "/resources", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	if key != "" {
		req.Header.Set("Idempotency-Key", key)
	}
	resp, err := app.Test(req)
	if err != nil {
		return nil, err
	}
	rec := httptest.NewRecorder()
	rec.WriteHeader(resp.StatusCode)
	for k, vs := range resp.Header {
		for _, v := range vs {
			rec.Header().Set(k, v)
		}
	}
	body, _ := io.ReadAll(resp.Body)
	rec.Write(body)
	return rec, nil
}

// ── No key ────────────────────────────────────────────────────────────────────

func TestIdempotency_NoKey_PassesThrough(t *testing.T) {
	store := idempotency.NewInMemoryStore(24 * time.Hour)
	app := newTestApp(store)

	resp, err := app.Test(httptest.NewRequest("POST", "/resources", nil))

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	assert.Empty(t, resp.Header.Get("X-Idempotency-Replayed"))
}

// ── First request is processed normally ───────────────────────────────────────

func TestIdempotency_FirstRequest_Processed(t *testing.T) {
	store := idempotency.NewInMemoryStore(24 * time.Hour)
	app := newTestApp(store)

	req := httptest.NewRequest("POST", "/resources", nil)
	req.Header.Set("Idempotency-Key", "uuid-111")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	assert.Empty(t, resp.Header.Get("X-Idempotency-Replayed"))
}

// ── Second request returns cached response ────────────────────────────────────

func TestIdempotency_SecondRequest_ReturnsCached(t *testing.T) {
	store := idempotency.NewInMemoryStore(24 * time.Hour)
	app := newTestApp(store)

	makeReq := func() *http.Request {
		r := httptest.NewRequest("POST", "/resources", nil)
		r.Header.Set("Idempotency-Key", "uuid-222")
		return r
	}

	// First call
	resp1, err := app.Test(makeReq())
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp1.StatusCode)

	// Second call — same key
	resp2, err := app.Test(makeReq())
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp2.StatusCode)
	assert.Equal(t, "true", resp2.Header.Get("X-Idempotency-Replayed"))
	assert.Equal(t, "uuid-222", resp2.Header.Get("Idempotency-Key"))

	body2, _ := io.ReadAll(resp2.Body)
	assert.Contains(t, string(body2), "created")
}

// ── Idempotency does NOT apply to GET / DELETE ─────────────────────────────────

func TestIdempotency_GET_Bypassed(t *testing.T) {
	store := idempotency.NewInMemoryStore(24 * time.Hour)
	app := newTestApp(store)

	makeReq := func() *http.Request {
		r := httptest.NewRequest("GET", "/resources", nil)
		r.Header.Set("Idempotency-Key", "uuid-get")
		return r
	}

	resp1, _ := app.Test(makeReq())
	resp2, _ := app.Test(makeReq())

	// Both should be live responses, never replayed
	assert.Equal(t, fiber.StatusOK, resp1.StatusCode)
	assert.Equal(t, fiber.StatusOK, resp2.StatusCode)
	assert.Empty(t, resp2.Header.Get("X-Idempotency-Replayed"))
}

func TestIdempotency_DELETE_Bypassed(t *testing.T) {
	store := idempotency.NewInMemoryStore(24 * time.Hour)
	app := newTestApp(store)

	makeReq := func() *http.Request {
		r := httptest.NewRequest("DELETE", "/resources/1", nil)
		r.Header.Set("Idempotency-Key", "uuid-del")
		return r
	}

	resp1, _ := app.Test(makeReq())
	resp2, _ := app.Test(makeReq())

	assert.Equal(t, fiber.StatusNoContent, resp1.StatusCode)
	assert.Equal(t, fiber.StatusNoContent, resp2.StatusCode)
	assert.Empty(t, resp2.Header.Get("X-Idempotency-Replayed"))
}

// ── Same key on different paths are independent ───────────────────────────────

func TestIdempotency_SameKey_DifferentPaths_Independent(t *testing.T) {
	store := idempotency.NewInMemoryStore(24 * time.Hour)
	app := newTestApp(store)

	postReq := httptest.NewRequest("POST", "/resources", nil)
	postReq.Header.Set("Idempotency-Key", "shared-uuid")

	putReq := httptest.NewRequest("PUT", "/resources/1", nil)
	putReq.Header.Set("Idempotency-Key", "shared-uuid")

	respPost, _ := app.Test(postReq)
	respPut, _ := app.Test(putReq)

	// Neither should be marked as a replay — they're different scoped keys
	assert.Empty(t, respPost.Header.Get("X-Idempotency-Replayed"))
	assert.Empty(t, respPut.Header.Get("X-Idempotency-Replayed"))
}

// ── Expired entries are reprocessed ──────────────────────────────────────────

func TestIdempotency_ExpiredEntry_Reprocessed(t *testing.T) {
	store := idempotency.NewInMemoryStore(1 * time.Millisecond)
	app := newTestApp(store)

	makeReq := func() *http.Request {
		r := httptest.NewRequest("POST", "/resources", nil)
		r.Header.Set("Idempotency-Key", "uuid-expire")
		return r
	}

	app.Test(makeReq()) // prime the cache

	time.Sleep(5 * time.Millisecond) // let the entry expire

	resp, _ := app.Test(makeReq())

	// After expiry the request is processed fresh — no replay header
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	assert.Empty(t, resp.Header.Get("X-Idempotency-Replayed"))
}

// ── Store errors are surfaced ─────────────────────────────────────────────────

func TestIdempotency_StoreGetError_Surfaced(t *testing.T) {
	store := idempotency.NewInMemoryStore(24 * time.Hour)
	store.OnGet = func(_ context.Context, _ string) (*idempotency.Entry, error) {
		return nil, fmt.Errorf("redis unavailable")
	}
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Use(Idempotency(store))
	app.Post("/resources", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusCreated)
	})

	req := httptest.NewRequest("POST", "/resources", nil)
	req.Header.Set("Idempotency-Key", "uuid-err")
	resp, err := app.Test(req)

	assert.NoError(t, err) // app.Test itself doesn't fail
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

// ── PUT is also covered ───────────────────────────────────────────────────────

func TestIdempotency_PUT_Cached(t *testing.T) {
	store := idempotency.NewInMemoryStore(24 * time.Hour)
	app := newTestApp(store)

	makeReq := func() *http.Request {
		r := httptest.NewRequest("PUT", "/resources/42", nil)
		r.Header.Set("Idempotency-Key", "uuid-put")
		return r
	}

	resp1, _ := app.Test(makeReq())
	resp2, _ := app.Test(makeReq())

	assert.Equal(t, fiber.StatusOK, resp1.StatusCode)
	assert.Equal(t, fiber.StatusOK, resp2.StatusCode)
	assert.Equal(t, "true", resp2.Header.Get("X-Idempotency-Replayed"))
}
