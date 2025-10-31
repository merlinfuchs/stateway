package broker

type ServiceType string

const (
	SerivceTypeGateway ServiceType = "gateway"
	ServiceTypeCache   ServiceType = "cache"
)
