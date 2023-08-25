package main

import (
    "context"
    "log"

    "github.com/ferchox920/ecommerce-go/database"
    "github.com/ferchox920/ecommerce-go/routes"
    "github.com/ferchox920/ecommerce-go/service"
    "github.com/gofiber/fiber/v2"
)

func main() {
    app := fiber.New()

    client, err := database.ConnectDatabase()
    if err != nil {
        log.Fatalf("Error connecting to the database: %v", err)
    }
    defer client.Disconnect(context.Background())

    userService := services.NewUserService(client)
    authService := services.NewAuthService(client) // Nuevo servicio de autenticación

    // Middleware de registro
    app.Use(func(c *fiber.Ctx) error {
        log.Printf("Request: %s %s", c.Method(), c.Path())
        return c.Next()
    })

    routes.SetupRoutes(app, userService) // Llamada a SetupRoutes con mayúscula
    routes.SetupAuthRoutes(app, authService) // Agregar rutas de autenticación

    app.Listen(":3000")
}
