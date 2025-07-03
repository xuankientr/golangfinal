package http

import (
	"strconv"

	"github.com/Hiendang123/golang-server.git/internal/common"
	"github.com/Hiendang123/golang-server.git/internal/domain"
	"github.com/Hiendang123/golang-server.git/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type TaskHandler struct {
	Usecase *usecase.TaskUsecase
}

func NewTaskHandler(app *fiber.App, uc *usecase.TaskUsecase) {
	handler := &TaskHandler{Usecase: uc}

	app.Post("/tasks", common.AuthMiddleware, handler.Create)
	app.Get("/tasks", common.AuthMiddleware, handler.GetAll)
	app.Get("/tasks/:id", common.AuthMiddleware, handler.GetByID)
	app.Put("/tasks/:id", common.AuthMiddleware, handler.Update)
	// app.Delete("/tasks/:id", handler.Delete)
	app.Delete("/tasks/all", common.AuthMiddleware, handler.DeleteAll)
}

// Create godoc
// @Summary Create a new task
// @Description Create a new task for the authenticated user
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body domain.Task true "Task object"
// @Success 201 {object} map[string]interface{} "Created task ID"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 422 {object} map[string]interface{} "Validation Error"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security BearerAuth
// @Router /tasks [post]
func (h *TaskHandler) Create(c *fiber.Ctx) error {
	var task domain.Task
	if err := c.BodyParser(&task); err != nil {
		return common.RespondError(c, fiber.StatusBadRequest, "Invalid JSON")
	}

	if err := task.Validate(); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	userID := c.Locals("user_id").(uint)
	task.CreatedBy = userID
	task.UpdatedBy = userID
	if err := h.Usecase.CreateTask(&task); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create task")
	}

	return common.ResponseCreate(c, task.ID)
}

func (h *TaskHandler) GetAll(c *fiber.Ctx) error {
	titleParam := c.Query("title")
	doneParam := c.Query("done")
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if limit <= 0 || limit > 100 {
		limit = 10
	}

	tasks, total, err := h.Usecase.GetAll(limit, offset, domain.Task{Title: titleParam, Done: doneParam == "true"})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(common.PaginationResponse{Data: tasks, Total: total, Limit: limit, Offset: offset})
}

func (h *TaskHandler) GetByID(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	task, err := h.Usecase.GetTaskByID(uint(id))

	if err != nil {
		return fiber.NewError(fiber.StatusAccepted, err.Error())
	}

	return common.ResponseSuccess(c, task)
}

func (h *TaskHandler) Update(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	task, err := h.Usecase.GetTaskByID(uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	if err := c.BodyParser(task); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := h.Usecase.UpdateTask(task); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return common.ResponseSuccess(c, task)
}

// func (h *TaskHandler) Delete(c *fiber.Ctx) error {
// 	id, _ := strconv.Atoi(c.Params("id"))
// 	if err := h.Usecase.DeleteTask(uint(id)); err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Task not found"})
// 	}
// 	return c.SendStatus(fiber.StatusNoContent)
// }

func (h *TaskHandler) DeleteAll(c *fiber.Ctx) error {

	if err := h.Usecase.DeleteAll(); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return common.ResponseContent(c)
}
