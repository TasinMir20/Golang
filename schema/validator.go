package schema

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Global validator instance
var validate = validator.New()

// ValidateStruct takes any struct and validates it based on its tags.
// If validation fails, it returns a map of field-to-error-message.
// If validation succeeds, it returns nil.
func ValidateStruct(payload interface{}) map[string]string {
	// Create a map to hold our organized errors
	errors := make(map[string]string)

	// Validate the struct
	err := validate.Struct(payload)
	if err != nil {
		// If the error is of type ValidationErrors, we can process it
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldErr := range validationErrors {
				field := strings.ToLower(fieldErr.Field())
				tag := fieldErr.Tag()
				param := fieldErr.Param()

				// Create a user-friendly message
				message := formatErrorMessage(tag, param)
				errors[field] = message
			}
		}
		return errors
	}

	// Return nil if there are no errors
	return nil
}

// formatErrorMessage creates a user-friendly message for a given validation tag.
func formatErrorMessage(tag, param string) string {
	switch tag {
	case "required":
		return "This field is required."
	case "email":
		return "This field must be a valid email address."
	case "min":
		return fmt.Sprintf("This field must be at least %s characters long.", param)
	case "max":
		return fmt.Sprintf("This field must not be more than %s characters long.", param)
	default:
		return fmt.Sprintf("Validation failed on the '%s' rule.", tag)
	}
}
