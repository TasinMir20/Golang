package routes

import (
	authRoutes "new-project-go/api/routes/auth"

	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(route fiber.Router) {
	route.All("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "API root!",
		})
	})

	auth := route.Group("/auth")
	authRoutes.SetupAuthRoutes(auth)

}
