/*
	package main

import (

	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

)

var db *gorm.DB

// User struct

	type User struct {
		ID        uint `gorm:"primaryKey"`
		Name      string
		Email     string `gorm:"uniqueIndex"`
		Password  string
		CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
		UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	}

// In-memory storage
var tasks = map[int]Task{}
var nextID = 1

func main() {

		var err error
		// Connect to the database
		dsn := "host=localhost user=postgres password=password dbname=postgres port=5432 sslmode=disable"
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("failed to connect to database:", err)
		}
		log.Println("Connected to database successfully")
		// Migrate the schema
		db.AutoMigrate(&User{})
		db.AutoMigrate(&Task{})

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
			if err := db.Create(&task).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create task"})
			}
			return c.Status(fiber.StatusCreated).JSON(task)
		})

		// Read All
		app.Get("/tasks", func(c *fiber.Ctx) error {
			var tasks []Task
			if err := db.Find(&tasks).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch tasks"})
			}
			return c.JSON(tasks)
		})

		// Read One
		app.Get("/tasks/:id", func(c *fiber.Ctx) error {
			id := c.Params("id")

			var task Task
			if err := db.First(&task, id).Error; err != nil {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
			}
			return c.JSON(task)
		})

		// Update
		app.Put("/tasks/:id", func(c *fiber.Ctx) error {
			id := c.Params("id")
			var task Task
			if err := db.First(&task, id).Error; err != nil {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
			}
			var update Task
			if err := c.BodyParser(&update); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON!"})
			}
			task.Title = update.Title
			task.Done = update.Done
			if err := db.Save(&task).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update task"})
			}
			return c.JSON(task)
		})

		// Delete
		app.Delete("/tasks/:id", func(c *fiber.Ctx) error {
			id := c.Params("id")
			if err := db.Delete(&Task{}, id).Error; err != nil {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
			}
			return c.SendStatus(fiber.StatusNoContent)
		})

		app.Listen(":3000")
	}
*/
package main

import (
	"github.com/Hiendang123/golang-server.git/internal/common"
	httpapp "github.com/Hiendang123/golang-server.git/internal/delivery/http"
	"github.com/Hiendang123/golang-server.git/internal/repository/postgres"
	"github.com/Hiendang123/golang-server.git/internal/usecase"
	"github.com/Hiendang123/golang-server.git/pkg/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: common.ErrorHandler,
	})
	app.Use(recover.New())

	app.Use(common.Logger)

	db := database.InitDB()
	userRepo := postgres.NewUserPostgresRepo(db)
	userUC := usecase.NewUserUsecase(userRepo)
	httpapp.NewUserHandler(app, userUC)

	taskRepo := postgres.NewTaskPostgresRepo(db)
	taskUC := usecase.NewTaskUsecase(taskRepo)
	httpapp.NewTaskHandler(app, taskUC)

	app.Listen(":3000")
}
