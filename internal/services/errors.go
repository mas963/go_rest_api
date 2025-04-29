package services

type ErrorResponse struct {
	Code string `json:"code"`
	Message string `json:"message"`
}

type AppError struct {
	Code string
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

var (
	ErrNotFound = &AppError{
		Code: "NOT_FOUND",
		Message: "Resource not Found",
	}
	ErrValidation = &AppError{
		Code: "VALIDATION_ERROR",
		Message: "Validation failed",
	}
	ErrUnauthorized = &AppError{
		Code: "UNAUTHORIZED",
		Message: "Invalid credentials",
	}
	ErrInternal = &AppError{
		Code: "INTERNAL_ERROR",
		Message: "Internal server error",
	}
)

func GetErrorCode(err error) string {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code
	}
	return ErrInternal.Code
}