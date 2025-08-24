package router

import (
	"new-project-go/api/routes"

	"github.com/gofiber/fiber/v3"
)

func SetupRouter(app *fiber.App) {
	app.All("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Server is running!",
		})
	})

	api := app.Group("/api")
	routes.SetupRoutes(api)

}
