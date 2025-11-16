package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ErrorCode string

const (
	ErrorCodeUnknown        ErrorCode = "unknown"
	ErrorCodeInternal       ErrorCode = "internal"
	ErrorCodeNotFound       ErrorCode = "not_found"
	ErrorCodeAlreadyExists  ErrorCode = "already_exists"
	ErrorCodeInvalidRequest ErrorCode = "invalid_request"
)

type Error struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details any       `json:"details,omitempty"`
}

func (e *Error) UnmarshalJSON(data []byte) error {
	var aux struct {
		Code    ErrorCode       `json:"code"`
		Message string          `json:"message"`
		Details json.RawMessage `json:"details,omitempty"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	e.Code = aux.Code
	e.Message = aux.Message

	if e.Code == ErrorCodeInvalidRequest && aux.Details != nil {
		var rawValidationErrors map[string]string
		if err := json.Unmarshal(aux.Details, &rawValidationErrors); err != nil {
			return err
		}

		validationErrors := make(ValidationErrors, len(rawValidationErrors))
		for field, message := range rawValidationErrors {
			validationErrors[field] = errors.New(message)
		}

		e.Details = validationErrors
	}

	return nil
}

func NewError(code ErrorCode, message string, details any) *Error {
	return &Error{Code: code, Message: message, Details: details}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func ErrUnknown(message string) *Error {
	return NewError(ErrorCodeUnknown, message, nil)
}

func ErrNotFound(message string) *Error {
	return NewError(ErrorCodeNotFound, message, nil)
}

func ErrInvalidRequest(message string, details error) *Error {
	return NewError(ErrorCodeInvalidRequest, message, details)
}

func GetErrorCode(err error) ErrorCode {
	var sErr *Error
	if errors.As(err, &sErr) {
		return sErr.Code
	}
	return ErrorCodeUnknown
}

func IsErrorCode(err error, codes ...ErrorCode) bool {
	var sErr *Error
	if errors.As(err, &sErr) {
		return slices.Contains(codes, sErr.Code)
	}
	return false
}

type ValidationErrors = validation.Errors
