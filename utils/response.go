package utils

import (
	"github.com/gofiber/fiber/v3"
)

type Response struct {
	Data       interface{} `json:"data"`
	Message    interface{} `json:"message"`
	Status     bool        `json:"status"`
	StatusCode int         `json:"status_code"`
	Details    interface{} `json:"details"`
}

func toBool(v any) bool {
	switch v := v.(type) {
	case bool:
		return v
	case int, int8, int16, int32, int64:
		return v != 0
	case float32, float64:
		return v != 0.0
	case string:
		return v != ""
	case nil:
		return false
	default:
		return true
	}
}

func SendResponse(c fiber.Ctx, data interface{}, message interface{}, status bool, statusCode int, details interface{}) error {

	if !toBool(message) {
		message = "Success"
	}

	if !toBool(statusCode) {
		statusCode = 200
	}

	if status != false {
		status = true
	}

	return c.Status(statusCode).JSON(Response{
		Data:       data,
		Message:    message,
		Status:     status,
		StatusCode: statusCode,
		Details:    details,
	})
}
