package task

import (
	"bufio"
	"encoding/json"

	"github.com/fabienbellanger/fiber-boilerplate/db"
	"github.com/fabienbellanger/fiber-boilerplate/entities"
	"github.com/fabienbellanger/fiber-boilerplate/stores"
	"github.com/fabienbellanger/fiber-boilerplate/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type TaskHandler struct {
	router fiber.Router
	store  stores.TaskStorer
	db     *db.DB
	logger *zap.Logger
}

// New returns a new TaskHandler
func New(r fiber.Router, task stores.TaskStorer, db *db.DB, logger *zap.Logger) TaskHandler {
	return TaskHandler{
		router: r,
		store:  task,
		db:     db,
		logger: logger,
	}
}

// Routes adds tasks routes
func (t *TaskHandler) Routes() {
	t.router.Post("", t.create())
	t.router.Get("", t.getAll())
	t.router.Get("/stream", t.getAllStream(t.db))
}

// create creates a new task.
func (t *TaskHandler) create() fiber.Handler {
	return func(c *fiber.Ctx) error {
		task := new(entities.TaskForm)
		if err := c.BodyParser(task); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Bad Request",
			})
		}

		errors := utils.ValidateStruct(*task)
		if errors != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Bad Request",
				Details: errors,
			})
		}

		// Data validation
		// ---------------
		if task.Name == "" {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Name cannot be empty",
			})
		}

		// Database insertion
		// ------------------
		newTask := entities.Task{
			Name:        task.Name,
			Description: task.Description,
		}

		if err := t.store.Create(&newTask); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error during task creation")
		}
		return c.JSON(newTask)
	}
}

// getAll lists all tasks.
func (t *TaskHandler) getAll() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tasks, err := t.store.ListAll()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error during tasks list")
		}

		return c.JSON(tasks)
	}
}

// getAllStream lists all tasks with a stream.
func (t *TaskHandler) getAllStream(db *db.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rows, err := t.store.ListAllRows()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error during tasks list with stream")
		}

		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
			w.WriteString("[")
			enc := json.NewEncoder(w)

			for i := 0; rows.Next(); i++ {
				if i > 0 {
					w.WriteString(",")
				}

				var task entities.Task
				if err := db.ScanRows(rows, &task); err != nil {
					continue
				}
				if err := enc.Encode(task); err != nil {
					continue
				}
			}
			w.WriteString("]")

			defer rows.Close()
		})

		return nil
	}
}
