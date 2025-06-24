package config

import (
	"github.com/spf13/viper"
	"log"
)

type option struct {
	configFolder []string
	configFile   string
	configType   string
}

type Option func(*option)

// LoadConfig load config
func LoadConfig(opts ...Option) Config {
	var cfg Config
	opt := &option{
		configFolder: getDefaultConfigFolder(),
		configFile:   getDefaultConfigFile(),
		configType:   getDefaultConfigType(),
	}

	for _, optFunc := range opts {
		optFunc(opt)
	}

	for _, folder := range opt.configFolder {
		viper.AddConfigPath(folder)
	}

	viper.SetConfigName(opt.configFile)
	viper.SetConfigType(opt.configType)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("failed to load config: %s", err)
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatalf("failed to unmarshal config: %s", err)
	}

	return cfg
}

func getDefaultConfigFolder() []string {
	return []string{"./files/config"}
}

func getDefaultConfigFile() string {
	return "config"
}

func getDefaultConfigType() string {
	return "yaml"
}

func WithConfigFolder(folder []string) Option {
	return func(opt *option) {
		opt.configFolder = folder
	}
}

func WithConfigFile(file string) Option {
	return func(opt *option) {
		opt.configFile = file
	}
}

func WithConfigType(configType string) Option {
	return func(opt *option) {
		opt.configType = configType
	}
}
