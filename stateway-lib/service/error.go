package service

import (
	"errors"
	"fmt"
	"slices"
)

type ErrorCode string

const (
	ErrorCodeUnknown       ErrorCode = "unknown"
	ErrorCodeInternal      ErrorCode = "internal"
	ErrorCodeNotFound      ErrorCode = "not_found"
	ErrorCodeAlreadyExists ErrorCode = "already_exists"
)

type Error struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func NewError(code ErrorCode, message string) *Error {
	return &Error{Code: code, Message: message}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func ErrNotFound(message string) *Error {
	return NewError(ErrorCodeNotFound, message)
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
