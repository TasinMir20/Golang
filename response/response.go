package response

import (
	"github.com/gofiber/fiber/v3"
)

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

func SendResponse(c fiber.Ctx, data interface{}, message interface{}, statusCode interface{}, status interface{}, errors interface{}) error {

	if !toBool(message) {
		message = "Success"
	}

	var status_code int
	if v, ok := statusCode.(int); ok {
		status_code = v
	} else {
		status_code = 200
	}

	if status != false {
		status = true
	}

	return c.Status(status_code).JSON(fiber.Map{
		"data":    data,
		"message": message,
		"status":  status,
		"errors":  errors,
	})

}
