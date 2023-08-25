// auth.go

package routes

import (
	services "github.com/ferchox920/ecommerce-go/services"
	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App, authService *services.AuthService) {
    app.Post("/auth/login", func(c *fiber.Ctx) error {
        var loginRequest struct {
            Email    string `json:"email"`
            Password string `json:"password"`
        }
        if err := c.BodyParser(&loginRequest); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "message": "Bad request",
            })
        }

        token, err := authService.Login(loginRequest.Email, loginRequest.Password)
        if err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "message": "Invalid credentials",
            })
        }

        return c.JSON(fiber.Map{
            "token": token,
        })
    })
}
