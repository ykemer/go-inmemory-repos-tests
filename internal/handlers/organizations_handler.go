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

func (h *OrganizationHandler) List(c *fiber.Ctx) error {
	res, err := h.service.List(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(res)
}

func (h *OrganizationHandler) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	res, err := h.service.Get(c.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Organization not found")
	}
	return c.JSON(res)
}

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
