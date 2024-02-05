package task

import (
	"bufio"
	"encoding/json"

	"github.com/fabienbellanger/fiber-boilerplate/domain/entities"
	"github.com/fabienbellanger/fiber-boilerplate/stores"
	"github.com/fabienbellanger/fiber-boilerplate/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type TaskHandler struct {
	router fiber.Router
	store  stores.TaskStorer
	logger *zap.Logger
}

// New returns a new TaskHandler
func New(r fiber.Router, task stores.TaskStorer, logger *zap.Logger) TaskHandler {
	return TaskHandler{
		router: r,
		store:  task,
		logger: logger,
	}
}

// Routes adds tasks routes
func (t *TaskHandler) Routes() {
	t.router.Post("", t.create())
	t.router.Get("", t.getAll())
	t.router.Get("/stream", t.getAllStream())
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
			return utils.NewError(c, t.logger, "Database error", "Error during task creation", err)
		}
		return c.JSON(newTask)
	}
}

// getAll lists all tasks.
func (t *TaskHandler) getAll() fiber.Handler {
	return func(c *fiber.Ctx) error {
		page := c.Query("p")
		limit := c.Query("l")
		sorts := c.Query("s")

		tasks, total, err := t.store.ListAll(page, limit, sorts)
		if err != nil {
			return utils.NewError(c, t.logger, "Database error", "Error during tasks list", err)
		}

		return c.JSON(utils.PaginateResponse{
			Total: total,
			Data:  tasks,
		})
	}
}

// getAllStream lists all tasks with a stream.
func (t *TaskHandler) getAllStream() fiber.Handler {
	return func(c *fiber.Ctx) error {
		rows, err := t.store.ListAllRows()
		if err != nil {
			return utils.NewError(c, t.logger, "Database error", "Error during tasks list with stream", err)
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
				if err := t.store.ScanRow(rows, &task); err != nil {
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
