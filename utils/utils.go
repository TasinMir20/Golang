package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"new-project-go/config"

	"new-project-go/models"
)

func LoginSessionCreate(userId primitive.ObjectID, expireInDays int) (string, error) {

	if expireInDays == 0 {
		expireInDays = 30
	}

	var UserSessionCollection = config.Database.Collection("usersessions")

	sessionUUID := uuid.NewString()
	expireDate := time.Now().AddDate(0, 0, expireInDays)

	session := models.UserSession{
		User:        userId,
		SessionName: "LoginSession",
		UUID:        sessionUUID,
		ExpireDate:  expireDate,
	}

	result, err := UserSessionCollection.InsertOne(context.Background(), session)
	if err != nil {
		return "", err
	}

	sessionID := result.InsertedID.(primitive.ObjectID).Hex()

	token, err := CreateSecretToken(sessionID, sessionUUID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func ParseBody[T any](c fiber.Ctx, payload *T) error {
	if err := c.Bind().Body(payload); err != nil {
		var unmarshalErr *json.UnmarshalTypeError
		if errors.As(err, &unmarshalErr) {
			return fiber.NewError(
				fiber.StatusBadRequest,
				fmt.Sprintf("Invalid type for field '%s'. Expected a %s but got a %s.",
					unmarshalErr.Field,
					unmarshalErr.Type.Name(),
					unmarshalErr.Value,
				),
			)
		}
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	return nil
}
