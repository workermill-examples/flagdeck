package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// ErrorCode represents standardized API error codes
type ErrorCode string

const (
	ErrorCodeValidation    ErrorCode = "VALIDATION_ERROR"
	ErrorCodeNotFound      ErrorCode = "NOT_FOUND"
	ErrorCodeUnauthorized  ErrorCode = "UNAUTHORIZED"
	ErrorCodeForbidden     ErrorCode = "FORBIDDEN"
	ErrorCodeConflict      ErrorCode = "CONFLICT"
	ErrorCodeRateLimit     ErrorCode = "RATE_LIMIT_EXCEEDED"
	ErrorCodeInternal      ErrorCode = "INTERNAL_ERROR"
	ErrorCodeDatabaseError ErrorCode = "DATABASE_ERROR"
	ErrorCodeBadRequest    ErrorCode = "BAD_REQUEST"
)

// ErrorDetail represents the standardized error response format
type ErrorDetail struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

// ErrorResponse represents the standardized error response wrapper
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// CreateErrorResponse creates a standardized error response
func CreateErrorResponse(code ErrorCode, message string) ErrorResponse {
	return ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	}
}

// SendValidationError sends a validation error response
func SendValidationError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(
		CreateErrorResponse(ErrorCodeValidation, message),
	)
}

// SendNotFoundError sends a not found error response
func SendNotFoundError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusNotFound).JSON(
		CreateErrorResponse(ErrorCodeNotFound, message),
	)
}

// SendUnauthorizedError sends an unauthorized error response
func SendUnauthorizedError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(
		CreateErrorResponse(ErrorCodeUnauthorized, message),
	)
}

// SendForbiddenError sends a forbidden error response
func SendForbiddenError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusForbidden).JSON(
		CreateErrorResponse(ErrorCodeForbidden, message),
	)
}

// SendConflictError sends a conflict error response
func SendConflictError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusConflict).JSON(
		CreateErrorResponse(ErrorCodeConflict, message),
	)
}

// SendRateLimitError sends a rate limit exceeded error response
func SendRateLimitError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusTooManyRequests).JSON(
		CreateErrorResponse(ErrorCodeRateLimit, message),
	)
}

// SendInternalError sends an internal server error response
func SendInternalError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusInternalServerError).JSON(
		CreateErrorResponse(ErrorCodeInternal, message),
	)
}

// SendDatabaseError sends a database error response
func SendDatabaseError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusInternalServerError).JSON(
		CreateErrorResponse(ErrorCodeDatabaseError, message),
	)
}

// SendBadRequestError sends a bad request error response
func SendBadRequestError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(
		CreateErrorResponse(ErrorCodeBadRequest, message),
	)
}
