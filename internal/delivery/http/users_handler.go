package http

import (
	"github.com/Hiendang123/golang-server.git/internal/common"
	"github.com/Hiendang123/golang-server.git/internal/domain"
	"github.com/Hiendang123/golang-server.git/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	Usecase *usecase.UserUsecase
}

func NewUserHandler(app *fiber.App, uc *usecase.UserUsecase) {
	handler := &UserHandler{Usecase: uc}

	app.Post("/users/register", handler.CreateUser)
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var user domain.User
	if err := c.BodyParser(&user); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON")
	}

	if err := user.Validate(); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	if err := h.Usecase.CreateUser(&user); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return common.RespondCreated(c, user.ID)
}
