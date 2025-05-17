package main

import "github.com/gofiber/fiber/v2"

// Task struct
type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

// In-memory storage
var tasks = map[int]Task{}
var nextID = 1

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// Create
	app.Post("/tasks", func(c *fiber.Ctx) error {
		var task Task
		if err := c.BodyParser(&task); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON!"})
		}
		task.ID = nextID
		nextID++
		tasks[task.ID] = task
		return c.Status(fiber.StatusCreated).JSON(task)
	})

	// Read All
	app.Get("/tasks", func(c *fiber.Ctx) error {
		taskList := []Task{}
		for _, t := range tasks {
			taskList = append(taskList, t)
		}
		return c.JSON(taskList)
	})

	// Read One
	app.Get("/tasks/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse ID"})
		}
		task, exists := tasks[id]
		if !exists {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
		}
		return c.JSON(task)
	})

	// Update
	app.Put("/tasks/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse ID"})
		}
		var update Task
		if err := c.BodyParser(&update); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		task, exists := tasks[id]
		if !exists {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
		}

		task.Title = update.Title
		task.Done = update.Done
		tasks[id] = task
		return c.JSON(task)
	})

	// Delete
	app.Delete("/tasks/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse ID"})
		}
		_, exists := tasks[id]
		if !exists {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
		}
		delete(tasks, id)
		return c.SendStatus(fiber.StatusNoContent)
	})

	app.Listen(":3000")
}
