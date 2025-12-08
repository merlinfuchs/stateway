package service

import (
	"encoding/json"
)

type ServiceType string

const (
	ServiceTypeGateway ServiceType = "gateway"
	ServiceTypeCache   ServiceType = "cache"
	ServiceTypeAudit   ServiceType = "audit"
)

type Response struct {
	Success bool            `json:"success"`
	Error   *Error          `json:"error,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
}

type RequestValidate interface {
	Validate() error
}
