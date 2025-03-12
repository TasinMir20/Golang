package main

import (
	"fiber-api/routes"

	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New()

	routes.SetupRoutes(app)

	if err := app.Listen(":4000"); err != nil {
		panic(err)
	}
}
