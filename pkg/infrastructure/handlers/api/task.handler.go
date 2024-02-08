package api

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/domain/entities"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/domain/requests"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/domain/usecases"

	"github.com/fabienbellanger/fiber-boilerplate/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// Task handler
type Task struct {
	router      fiber.Router
	taskUseCase usecases.Task
	logger      *zap.Logger
}

// NewTask returns a new Handler
func NewTask(r fiber.Router, taskUseCase usecases.Task, logger *zap.Logger) Task {
	return Task{
		router:      r,
		taskUseCase: taskUseCase,
		logger:      logger,
	}
}

// TaskProtectedRoutes adds tasks routes
func (t *Task) TaskProtectedRoutes() {
	t.router.Post("", t.create())
	t.router.Get("", t.getAll())
	t.router.Get("/stream", t.getAllStream())
}

// create creates a new task.
func (t *Task) create() fiber.Handler {
	return func(c *fiber.Ctx) error {
		task := new(requests.TaskCreation)
		if err := c.BodyParser(task); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Bad Request",
			})
		}

		newTask, err := t.taskUseCase.Create(*task)
		if err != nil {
			if errors.Is(err, utils.HTTPError{}) && err.Err != nil {
				if details, ok := err.Details.(string); ok {
					return utils.NewError(c, t.logger, err.Message, details, err.Err)
				}
			}
			return c.Status(err.Code).JSON(err)
		}

		return c.JSON(newTask)
	}
}

// getAll lists all tasks.
func (t *Task) getAll() fiber.Handler {
	return func(c *fiber.Ctx) error {
		pagination := new(requests.Pagination)
		if err := c.QueryParser(pagination); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Bad Request",
			})
		}

		res, err := t.taskUseCase.GetAll(*pagination)
		if err != nil {
			if errors.Is(err, utils.HTTPError{}) && err.Err != nil {
				if details, ok := err.Details.(string); ok {
					return utils.NewError(c, t.logger, err.Message, details, err.Err)
				}
			}
			return c.Status(err.Code).JSON(err)
		}

		return c.JSON(res)
	}
}

// getAllStream lists all tasks with a stream.
func (t *Task) getAllStream() fiber.Handler {
	return func(c *fiber.Ctx) error {
		rows, err := t.taskUseCase.GetAllStream()
		if err != nil {
			if errors.Is(err, utils.HTTPError{}) && err.Err != nil {
				if details, ok := err.Details.(string); ok {
					return utils.NewError(c, t.logger, err.Message, details, err.Err)
				}
			}
			return c.Status(err.Code).JSON(err)
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
				if err := t.taskUseCase.ScanTask(rows, &task); err != nil {
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
