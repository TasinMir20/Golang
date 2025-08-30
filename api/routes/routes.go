package routes

import (
	"fmt"
	authRoutes "new-project-go/api/routes/auth"

	"github.com/gofiber/fiber/v3"
)

var middleware = func(c fiber.Ctx) error {
	c.Set("X-Middleware", "Active")
	fmt.Println("Middleware")
	token := c.Get("Authorization")

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication required",
		})
	} else {
		return c.Next()
	}

}

func SetupRoutes(route fiber.Router) {
	route.All("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "API root!",
		})
	})

	auth := route.Group("/auth", middleware)
	authRoutes.SetupAuthRoutes(auth)

}
