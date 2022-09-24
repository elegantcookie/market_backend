package apperror

import (
	"encoding/json"
)

var (
	ErrInvalidTokenSignature = NewAppError(nil, "wrong token", "NS-000004", "token signature is invalid")
	ErrEmptyTokenString      = NewAppError(nil, "wrong token", "NS-000005", "token string is empty")
)

type AppError struct {
	Err              error  `json:"-"`
	Message          string `json:"message,omitempty"`
	DeveloperMessage string `json:"developer_message,omitempty"`
	Code             string `json:"code,omitempty"`
}

func NewAppError(err error, message, code, developerMessage string) *AppError {
	return &AppError{
		Err:              err,
		Code:             code,
		Message:          message,
		DeveloperMessage: developerMessage,
	}
}

func (e *AppError) Error() string {
	return e.Err.Error()
}

func (e *AppError) Unwrap() error { return e.Err }

func (e *AppError) Marshal() []byte {
	bytes, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return bytes
}

func BadRequestError(message string) *AppError {
	return NewAppError(nil, message, "NS-000002", "some thing wrong with auth data")
}

func systemError(developerMessage string) *AppError {
	return NewAppError(nil, "system error", "NS-000001", developerMessage)
}
