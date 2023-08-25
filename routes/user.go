package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/ferchox920/ecommerce-go/models"
	"github.com/ferchox920/ecommerce-go/service"
)

func SetupRoutes(app *fiber.App, userService *services.UserService) {
	app.Get("/users", func(c *fiber.Ctx) error {
		return c.SendString("Hola")
	})

	app.Post("/users", func(c *fiber.Ctx) error {
		var user models.User
		if err := c.BodyParser(&user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Bad request",
			})
		}

		if err := userService.CreateUser(&user); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error creating user",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(user)
	})
}
