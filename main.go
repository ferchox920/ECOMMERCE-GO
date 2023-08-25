package main

import (
	"context"
	"log"

	"github.com/ferchox920/ecommerce-go/database"
	"github.com/ferchox920/ecommerce-go/routes"
	"github.com/ferchox920/ecommerce-go/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.mongodb.org/mongo-driver/mongo"
)

type AppServices struct {
    UserService    *services.UserService
    AuthService    *services.AuthService
}

func initAppServices(client *mongo.Client) *AppServices {
    userService := services.NewUserService(client)
    authService := services.NewAuthService(client)
    return &AppServices{
        UserService: userService,
        AuthService: authService,
    }
}

func main() {
    app := fiber.New()

    client, err := database.ConnectDatabase()
    if err != nil {
        log.Fatalf("Error connecting to the database: %v", err)
    }
    defer client.Disconnect(context.Background())

    services := initAppServices(client)

	// Middleware de cors
	app.Use(cors.New(cors.Config{
        AllowOrigins: "*",
        AllowMethods: "GET,POST,PUT,DELETE",
        AllowHeaders: "Origin, Content-Type, Accept",
    }))

    // Middleware de registro
    app.Use(func(c *fiber.Ctx) error {
        log.Printf("Request: %s %s", c.Method(), c.Path())
        return c.Next()
    })

    routes.SetupRoutes(app, services.UserService)
    routes.SetupAuthRoutes(app, services.AuthService)

    app.Listen(":3000")
}
