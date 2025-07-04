package config

type Config struct {
	App            AppConfig      `mapstructure:"app" validate:"required"`
	Database       DatabaseConfig `mapstructure:"database" validate:"required"`
	Redis          RedisConfig    `mapstructure:"redis" validate:"required"`
	Secrete        SecretConfig   `mapstructure:"secrete" validate:"required"`
	ProductService ProductService `mapstructure:"product_service" validate:"required"`
}

type ProductService struct {
	Host string `mapstructure:"host" validate:"required"`
}

type AppConfig struct {
	Port string `mapstructure:"port" validate:"required"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host" validate:"required"`
	Port     string `mapstructure:"port" validate:"required"`
	Name     string `mapstructure:"name" validate:"required"`
	Password string `mapstructure:"password" validate:"required"`
	User     string `mapstructure:"user" validate:"required"`
}

type RedisConfig struct {
	Port     string `mapstructure:"port" validate:"required"`
	Host     string `mapstructure:"host" validate:"required"`
	Password string `mapstructure:"password" validate:"required"`
}

type SecretConfig struct {
	JWTSecret string `mapstructure:"jwtsecret" validate:"required"`
}
