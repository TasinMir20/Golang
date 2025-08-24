package main

import (
	"new-project-go/router"

	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New()

	router.SetupRouter(app)

	if err := app.Listen(":4000"); err != nil {
		panic(err)
	}
}
