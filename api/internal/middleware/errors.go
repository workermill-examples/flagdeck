package middleware

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
)

type ErrorCode string

const (
	ErrorCodeInternal      ErrorCode = "INTERNAL_ERROR"
	ErrorCodeValidation    ErrorCode = "VALIDATION_ERROR"
	ErrorCodeNotFound      ErrorCode = "NOT_FOUND"
	ErrorCodeUnauthorized  ErrorCode = "UNAUTHORIZED"
	ErrorCodeForbidden     ErrorCode = "FORBIDDEN"
	ErrorCodeBadRequest    ErrorCode = "BAD_REQUEST"
	ErrorCodeConflict      ErrorCode = "CONFLICT"
	ErrorCodeRateLimit     ErrorCode = "RATE_LIMIT_EXCEEDED"
	ErrorCodeDatabaseError ErrorCode = "DATABASE_ERROR"
)

type APIError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

type ErrorResponse struct {
	Error APIError `json:"error"`
}

type CustomError struct {
	Code       ErrorCode
	Message    string
	StatusCode int
}

func (e CustomError) Error() string {
	return e.Message
}

func NewCustomError(code ErrorCode, message string, statusCode int) *CustomError {
	return &CustomError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

func NewValidationError(message string) *CustomError {
	return NewCustomError(ErrorCodeValidation, message, fiber.StatusBadRequest)
}

func NewNotFoundError(message string) *CustomError {
	return NewCustomError(ErrorCodeNotFound, message, fiber.StatusNotFound)
}

func NewUnauthorizedError(message string) *CustomError {
	return NewCustomError(ErrorCodeUnauthorized, message, fiber.StatusUnauthorized)
}

func NewForbiddenError(message string) *CustomError {
	return NewCustomError(ErrorCodeForbidden, message, fiber.StatusForbidden)
}

func NewBadRequestError(message string) *CustomError {
	return NewCustomError(ErrorCodeBadRequest, message, fiber.StatusBadRequest)
}

func NewConflictError(message string) *CustomError {
	return NewCustomError(ErrorCodeConflict, message, fiber.StatusConflict)
}

func NewRateLimitError(message string) *CustomError {
	return NewCustomError(ErrorCodeRateLimit, message, fiber.StatusTooManyRequests)
}

func NewDatabaseError(message string) *CustomError {
	return NewCustomError(ErrorCodeDatabaseError, message, fiber.StatusInternalServerError)
}

func NewInternalError(message string) *CustomError {
	return NewCustomError(ErrorCodeInternal, message, fiber.StatusInternalServerError)
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	var customErr *CustomError
	var fiberErr *fiber.Error

	if errors.As(err, &customErr) {
		return c.Status(customErr.StatusCode).JSON(ErrorResponse{
			Error: APIError{
				Code:    customErr.Code,
				Message: customErr.Message,
			},
		})
	}

	if errors.As(err, &fiberErr) {
		var code ErrorCode
		switch fiberErr.Code {
		case fiber.StatusBadRequest:
			code = ErrorCodeBadRequest
		case fiber.StatusUnauthorized:
			code = ErrorCodeUnauthorized
		case fiber.StatusForbidden:
			code = ErrorCodeForbidden
		case fiber.StatusNotFound:
			code = ErrorCodeNotFound
		case fiber.StatusConflict:
			code = ErrorCodeConflict
		case fiber.StatusTooManyRequests:
			code = ErrorCodeRateLimit
		default:
			code = ErrorCodeInternal
		}

		return c.Status(fiberErr.Code).JSON(ErrorResponse{
			Error: APIError{
				Code:    code,
				Message: fiberErr.Message,
			},
		})
	}

	log.Printf("Unhandled error: %v", err)

	return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
		Error: APIError{
			Code:    ErrorCodeInternal,
			Message: "Internal server error",
		},
	})
}
