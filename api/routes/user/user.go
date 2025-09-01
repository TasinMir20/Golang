package userRoutes

import (
	userControllers "new-project-go/api/controllers/user"

	"github.com/gofiber/fiber/v3"
)

func SetupAuthRoutes(route fiber.Router) {

	route.Get("/", userControllers.GetUser)

}
