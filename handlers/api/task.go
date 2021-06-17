package api

import (
	"bufio"
	"encoding/json"

	"github.com/fabienbellanger/fiber-boilerplate/db"
	models "github.com/fabienbellanger/fiber-boilerplate/models"
	"github.com/fabienbellanger/fiber-boilerplate/repositories"
	"github.com/fabienbellanger/fiber-boilerplate/utils"
	"github.com/gofiber/fiber/v2"
)

// GetAllTasks lists all tasks.
// @Summary List all tasks
// @Description List all tasks
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {array} models.Task
// @Failure 400 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Security ApiKeyAuth
// @Router /tasks [get]
func GetAllTasks(db *db.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tasks, err := repositories.ListAllTasks(db)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error during tasks list")
		}

		return c.JSON(tasks)
	}
}

// GetAllTasksStream lists all tasks with a stream.
func GetAllTasksStream(db *db.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rows, err := repositories.ListAllTasksRows(db)
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

				var task models.Task
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

// CreateTask creates a new task.
func CreateTask(db *db.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		task := new(models.TaskForm)
		if err := c.BodyParser(task); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Bad Request",
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
		newTask := models.Task{
			Name:        task.Name,
			Description: task.Description,
		}

		if err := repositories.CreateTask(db, &newTask); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error during task creation")
		}
		return c.JSON(newTask)
	}
}
