package domain

// AppConfig - represents application configuration
type AppConfig struct {
	Environment   string `env:"ENVIRONMENT" validate:"required"`
	Port          string `env:"PORT" validate:"required"`
	DSN           string `env:"DSN" validate:"required"`
	RedisHost     string `env:"REDIS_HOST" validate:"required"`
	RedisPort     string `env:"REDIS_PORT" validate:"required"`
	RedisPassword string `env:"REDIS_PASSWORD" validate:"required"`
	RedisPoolSize int    `env:"REDIS_POOLSIZE" validate:"required"`
	RedisDB       int    `env:"REDIS_DB" validate:"required"`
	RabbitMQURI   string `env:"RABBITMQ_URI" validate:"required"`
}
