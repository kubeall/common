package messagebus

const (
	MessageBusDapr     = "dapr"
	MessageBusKafka    = "kafka"
	MessageBusPulsar   = "pulsar"
	MessageBusRocketMQ = "rocketmq"
	MessageBusRedis    = "redis"
	MessageBusNats     = "nats"
	MessageBusRabbitMQ = "rabbitmq"
)
const (
	TopicOrganization            = "/messagebus/eauth/organization"
	TopicOrganizationAccount     = "/messagebus/eauth/account"
	TopicOrganizationWorkspace   = "/messagebus/eauth/workspace"
	TopicOrganizationApplication = "/messagebus/eauth/application"
)
