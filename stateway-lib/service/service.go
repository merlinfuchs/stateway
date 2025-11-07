package service

type ServiceType string

const (
	ServiceTypeGateway ServiceType = "gateway"
	ServiceTypeCache   ServiceType = "cache"
)
