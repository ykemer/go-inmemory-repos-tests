package handlers

import (
	"tests/internal/dtos"
	"tests/internal/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

type OrganizationHandler struct {
	service *services.OrganizationService
}

func NewOrganizationHandler(service *services.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{service: service}
}

// List godoc
// @Summary      List all organizations
// @Tags         Organizations
// @Produce      json
// @Success      200  {array}   dtos.OrganizationResponse
// @Router       /organizations [get]
func (h *OrganizationHandler) List(c *fiber.Ctx) error {
	res, err := h.service.List(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(res)
}

// Get godoc
// @Summary      Get an organization
// @Tags         Organizations
// @Produce      json
// @Param        id   path      int  true  "Organization ID"
// @Success      200  {object}  dtos.OrganizationResponse
// @Failure      404  {object}  dtos.ProblemDetail
// @Router       /organizations/{id} [get]
func (h *OrganizationHandler) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	res, err := h.service.Get(c.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Organization not found")
	}
	return c.JSON(res)
}

// Create godoc
// @Summary      Create an organization
// @Tags         Organizations
// @Accept       json
// @Produce      json
// @Param        Idempotency-Key  header  string                           false  "Client-generated UUID. Repeated requests with the same key return the cached response for 1 hour."
// @Param        body             body    dtos.CreateOrganizationRequest   true   "Request body"
// @Success      201  {object}  dtos.OrganizationResponse
// @Failure      422  {object}  dtos.ProblemDetail
// @Router       /organizations [post]
func (h *OrganizationHandler) Create(c *fiber.Ctx) error {
	req := new(dtos.CreateOrganizationRequest)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Cannot parse JSON")
	}

	if err := validate.Struct(req); err != nil {
		return handleValidationError(c, err)
	}

	res, err := h.service.Create(c.Context(), *req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

// Update godoc
// @Summary      Update an organization
// @Tags         Organizations
// @Accept       json
// @Produce      json
// @Param        id               path    int                              true   "Organization ID"
// @Param        Idempotency-Key  header  string                           false  "Client-generated UUID. Repeated requests with the same key return the cached response for 1 hour."
// @Param        body             body    dtos.UpdateOrganizationRequest   true   "Request body"
// @Success      200  {object}  dtos.OrganizationResponse
// @Failure      404  {object}  dtos.ProblemDetail
// @Failure      422  {object}  dtos.ProblemDetail
// @Router       /organizations/{id} [put]
func (h *OrganizationHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	req := new(dtos.UpdateOrganizationRequest)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Cannot parse JSON")
	}

	if err := validate.Struct(req); err != nil {
		return handleValidationError(c, err)
	}

	res, err := h.service.Update(c.Context(), id, *req)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Organization not found")
	}

	return c.JSON(res)
}

// Delete godoc
// @Summary      Delete an organization
// @Tags         Organizations
// @Produce      json
// @Param        id   path  int  true  "Organization ID"
// @Success      204
// @Failure      404  {object}  dtos.ProblemDetail
// @Router       /organizations/{id} [delete]
func (h *OrganizationHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.Delete(c.Context(), id); err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Organization not found")
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func handleValidationError(c *fiber.Ctx, err error) error {
	fields := make(map[string]string)
	for _, e := range err.(validator.ValidationErrors) {
		fields[e.Field()] = e.Tag()
	}
	c.Locals("invalid_params", fields)
	return fiber.NewError(fiber.StatusUnprocessableEntity, "Validation failed")
}
