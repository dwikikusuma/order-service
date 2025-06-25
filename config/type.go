package config

type Config struct {
	App      AppConfig      `yaml:"app" validate:"required"`
	Database DatabaseConfig `yaml:"database" validate:"required"`
	Redis    RedisConfig    `yaml:"redis" validate:"required"`
	Secrete  SecretConfig   `yaml:"secrete" validate:"required"`
}

type AppConfig struct {
	Port string `yaml:"port" validate:"required"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host" validate:"required"`
	Port     string `yaml:"port" validate:"required"`
	Name     string `yaml:"name" validate:"required"`
	Password string `yaml:"password" validate:"required"`
	User     string `yaml:"user" validate:"required"`
}

type RedisConfig struct {
	Port     string `yaml:"port" validate:"required"`
	Host     string `yaml:"host" validate:"required"`
	Password string `yaml:"password" validate:"required"`
}

type SecretConfig struct {
	JWTSecret string `yaml:"jwtsecret" validate:"required"`
}
