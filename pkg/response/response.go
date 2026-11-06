package response

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var errorNames = map[int]string{
	http.StatusBadRequest:       "Bad Request",
	http.StatusUnauthorized:     "Unauthorized",
	http.StatusForbidden:        "Forbidden",
	http.StatusNotFound:         "Not Found",
	http.StatusMethodNotAllowed: "Method Not Allowed",
	http.StatusConflict:         "Conflict",
}

type Response struct {
	Success bool       `json:"success"`
	Data    any        `json:"data,omitempty"`
	Error   *ErrorInfo `json:"error,omitempty"`
}

type ErrorInfo struct {
	Name    string  `json:"name"`
	Message *string `json:"message,omitempty"`
}

func OK(c *gin.Context, statusCode int, data any) {
	c.JSON(statusCode, Response{
		Success: true,
		Data:    data,
	})
}

func Fail(c *gin.Context, statusCode int, err error) {
	errName := "Internal Server Error"
	var errMessage *string

	if value, ok := errorNames[statusCode]; ok {
		errName = value
	}

	if err != nil {
		message := err.Error()

		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			message = validationErrorsHandler(validationErrors)
		}

		errMessage = &message
	}

	c.JSON(statusCode, Response{
		Success: false,
		Error:   &ErrorInfo{Name: errName, Message: errMessage},
	})

	c.Abort()
}

func validationErrorsHandler(validationErrors validator.ValidationErrors) string {
	for _, e := range validationErrors {
		var fieldName = strings.ToLower(e.Field())

		switch e.Tag() {
		case "required":
			return "'" + fieldName + "' field is required"
		case "email":
			return "'" + fieldName + "' field should be an email"
		case "min":
			return "'" + fieldName + "' field length should be greater then or equal to " + e.Param()
		case "max":
			return "'" + fieldName + "' field length should be lower then or equal to " + e.Param()
		}
	}

	return validationErrors.Error()
}
