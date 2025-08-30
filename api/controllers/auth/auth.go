package authControllers

import (
	"context"
	"fmt"
	"new-project-go/config"
	"new-project-go/models"
	"new-project-go/response"
	"new-project-go/schema"
	"new-project-go/utils"
	"strings"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson"
)

func SignUp(c fiber.Ctx) error {
	var payload schema.SignInInput

	// Parse request body
	if err := utils.ParseBody(c, &payload); err != nil {
		return response.SendResponse(c, nil, err.Error(), 400, false, nil)
	}

	payload.Email = strings.ToLower(payload.Email)

	// Validate payload
	errors := schema.ValidateStruct(payload)
	if errors != nil {
		const message = "Invalid payload!"
		return response.SendResponse(c, nil, message, 400, false, errors)
	}

	var UserCollection = config.Database.Collection("users")

	// Check if user already exists
	existingUser := &models.User{}
	err := UserCollection.Find(context.Background(), bson.M{"email": payload.Email}).One(existingUser)
	if err == nil {
		const message = "User with this email already exists"
		return response.SendResponse(c, nil, message, 409, false, nil)
	}

	// Create new user
	user := &models.User{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: payload.Password,
	}

	result, err := UserCollection.InsertOne(context.Background(), user)
	if err != nil {
		fmt.Println(err)
		return response.SendResponse(c, nil, err.Error(), 500, nil, nil)
	}

	// Get the created user (without password)
	createdUser := &models.User{}
	UserCollection.Find(context.Background(), bson.M{"_id": result.InsertedID}).One(createdUser)

	const message = "User created successfully"
	return response.SendResponse(c, createdUser, message, 201, nil, nil)
}

func Login(c fiber.Ctx) error {
	var payload schema.LoginInput

	// Parse request body
	if err := utils.ParseBody(c, &payload); err != nil {
		return response.SendResponse(c, nil, err.Error(), 400, false, nil)
	}

	payload.Email = strings.ToLower(payload.Email)

	// Validate payload
	errors := schema.ValidateStruct(payload)
	if errors != nil {
		const message = "Invalid payload!"
		return response.SendResponse(c, nil, message, 400, false, errors)
	}

	var UserCollection = config.Database.Collection("users")

	// Find user by email
	user := &models.User{}
	err := UserCollection.Find(context.Background(), bson.M{"email": payload.Email}).One(user)
	if err != nil {
		return response.SendResponse(c, nil, "Invalid email or password", 401, false, nil)
	}

	// Check password
	if !user.CheckPassword(payload.Password) {
		return response.SendResponse(c, nil, "Password wrong!", 401, false, nil)
	}

	userId := user.ID
	token, err := utils.LoginSessionCreate(userId, 30)
	if err != nil {
		return response.SendResponse(c, nil, err, 500, false, nil)
	}

	var data = map[string]string{"token": token}

	return response.SendResponse(c, data, "Login successful!", nil, nil, nil)

}
