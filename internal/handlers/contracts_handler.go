package handlers

import (
	"tests/internal/dtos"
	"tests/internal/services"

	"github.com/gofiber/fiber/v2"
)

type ContractHandler struct {
	service *services.ContractService
}

func NewContractHandler(service *services.ContractService) *ContractHandler {
	return &ContractHandler{service: service}
}

func (h *ContractHandler) List(c *fiber.Ctx) error {
	orgId := c.Params("orgId")
	res, err := h.service.ListByOrg(c.Context(), orgId)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Organization not found")
	}
	return c.JSON(res)
}

func (h *ContractHandler) Get(c *fiber.Ctx) error {
	orgId := c.Params("orgId")
	contractId := c.Params("contractId")
	res, err := h.service.Get(c.Context(), orgId, contractId)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Contract not found")
	}
	return c.JSON(res)
}

func (h *ContractHandler) Create(c *fiber.Ctx) error {
	orgId := c.Params("orgId")
	req := new(dtos.CreateContractRequest)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Cannot parse JSON")
	}

	if err := validate.Struct(req); err != nil {
		return handleValidationError(c, err)
	}

	res, err := h.service.Create(c.Context(), orgId, *req)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

func (h *ContractHandler) Update(c *fiber.Ctx) error {
	orgId := c.Params("orgId")
	contractId := c.Params("contractId")
	req := new(dtos.UpdateContractRequest)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Cannot parse JSON")
	}

	if err := validate.Struct(req); err != nil {
		return handleValidationError(c, err)
	}

	res, err := h.service.Update(c.Context(), orgId, contractId, *req)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Contract not found")
	}

	return c.JSON(res)
}

func (h *ContractHandler) Delete(c *fiber.Ctx) error {
	orgId := c.Params("orgId")
	contractId := c.Params("contractId")
	if err := h.service.Delete(c.Context(), orgId, contractId); err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Contract not found")
	}
	return c.SendStatus(fiber.StatusNoContent)
}
