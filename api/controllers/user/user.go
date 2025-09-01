package userControllers

import (
	"new-project-go/models"
	"new-project-go/response"

	"github.com/gofiber/fiber/v3"
)

func GetUser(c fiber.Ctx) error {
	LoggedUser := c.Locals("LoggedUser").(*models.User)

	return response.SendResponse(c, LoggedUser, nil, nil, nil, nil)
}
