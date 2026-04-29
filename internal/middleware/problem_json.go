package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Problem struct {
	Type          string            `json:"type"`
	Title         string            `json:"title"`
	Status        int               `json:"status"`
	Detail        string            `json:"detail,omitempty"`
	InvalidParams map[string]string `json:"invalid_params,omitempty"`
}

// ProblemJSONErrorHandler is a Fiber ErrorHandler that formats every error
// as an RFC 7807 application/Problem+json response.
//
// Handlers signal validation failures by storing a map[string]string in
// c.Locals("invalid_params") before returning the error.
func ProblemJSONErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	detail := "An unexpected error occurred"

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		code = fiberErr.Code
		detail = fiberErr.Message
	}

	p := Problem{
		Type:   fmt.Sprintf("https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status/%d", code),
		Title:  http.StatusText(code),
		Status: code,
		Detail: detail,
	}

	if v := c.Locals("invalid_params"); v != nil {
		if fields, ok := v.(map[string]string); ok {
			p.InvalidParams = fields
		}
	}

	return c.Status(code).JSON(p, "application/Problem+json")
}
