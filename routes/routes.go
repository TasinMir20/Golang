package routes

import (
	"fiber-api/handlers"

	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", handlers.HandleHome)
	app.Get("/greet/:name", handlers.HandleGreet)
	app.Get("/fibonacci/:n", handlers.HandleFibonacci)
	app.Get("/stress/:seconds", handlers.HandleStress)
}
