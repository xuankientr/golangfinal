package common

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type PaginationResponse struct {
	Data   any   `json:"data"`
	Total  int64 `json:"total"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
}

func ResponseCreated(c *fiber.Ctx, data any) error {
	return c.Status(fiber.StatusCreated).JSON(data)
}

func ResponseSuccess(c *fiber.Ctx, data any) error {
	if data == nil {
		return c.Status(fiber.StatusNoContent).JSON(data)
	}
	return c.Status(fiber.StatusOK).JSON(data)
}

func ResponseContent(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusNoContent)
}
