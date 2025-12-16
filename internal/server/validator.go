package server

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterValidation("strongpassword", validateStrongPassword)
}

func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':\"\\|,.<>\/?]`).MatchString(password)

	return hasUpper && hasLower && hasNumber && hasSpecial
}

func (s *Server) validateStruct(v any) error {
	return validate.Struct(v)
}

func (s *Server) formatValidationErrors(err error) map[string]any {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		errors := make(map[string]map[string]string)
		for _, e := range validationErrors {
			field := strings.ToLower(e.Field())
			errors[field] = map[string]string{
				"key":   fmt.Sprintf("validation.%s", e.Tag()),
				"field": field,
				"param": e.Param(),
			}
		}
		return map[string]any{
			"statusCode": http.StatusBadRequest,
			"message":    "Validation failed",
			"errors":     errors,
		}
	}
	return map[string]any{
		"statusCode": http.StatusBadRequest,
		"message":    err.Error(),
	}
}
