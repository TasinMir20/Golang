package authRoutes

import (
	authControllers "new-project-go/api/controllers/auth"

	"github.com/gofiber/fiber/v3"
)

func SetupAuthRoutes(route fiber.Router) {

	route.Post("/login", authControllers.Login)
	route.Post("/sign-up", authControllers.SignUp)

}
