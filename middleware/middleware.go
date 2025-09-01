package middleware

import (
	"context"
	"log"
	"new-project-go/config"
	"new-project-go/models"
	"new-project-go/response"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"new-project-go/utils"
)

var AuthMiddleware = func(c fiber.Ctx) error {
	c.Set("X-Middleware", "Active")
	token := c.Get("Authorization")

	if token == "" {
		msg := "Please provide token in - headers.authorization"
		return response.SendResponse(c, nil, msg, 401, nil, nil)
	}

	JwtPayload, err := utils.ParseToken(token)
	if err != nil {
		log.Printf("Error: %v", err.Error())
		return response.SendResponse(c, nil, err.Error(), 401, nil, nil)

	}

	sessionId, exists := JwtPayload["sessionId"]
	if !exists {
		return response.SendResponse(c, nil, "sessionId not found in token", 401, nil, nil)
	}

	sessionIdStr, ok := sessionId.(string)
	if !ok {
		return response.SendResponse(c, nil, "sessionId is not a string", 401, false, nil)
	}

	sessionIdObjectID, err := primitive.ObjectIDFromHex(sessionIdStr)
	if err != nil {
		return response.SendResponse(c, nil, "Invalid sessionId format!", 401, false, nil)
	}

	var UserSessionCollection = config.Database.Collection("usersessions")

	pipeline := qmgo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"_id": sessionIdObjectID}}},
		bson.D{{Key: "$lookup", Value: bson.M{
			"from":         "users",
			"localField":   "user",
			"foreignField": "_id",
			"as":           "userData",
		}}},
		bson.D{{Key: "$unwind", Value: "$userData"}},
	}

	var userSessions []models.UserSession
	err = UserSessionCollection.Aggregate(context.Background(), pipeline).All(&userSessions)
	if err != nil {
		return response.SendResponse(c, nil, "Invalid token!", 401, false, nil)
	}

	if len(userSessions) == 0 {
		return response.SendResponse(c, nil, "Invalid token!", 401, false, nil)
	}

	userSession := userSessions[0]

	if userSession.ExpireDate.After(time.Now()) {
		c.Locals("LoggedUser", userSession.UserData)

		return c.Next()
	} else {
		return response.SendResponse(c, nil, "Session Expired!", 401, nil, nil)
	}

}
