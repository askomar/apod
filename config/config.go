package config

import (
	"github.com/caarlos0/env/v11"
)

type (
	Config struct {
		Database                 DatabaseConfig `envPrefix:"DB_"`
		Minio                    MinioConfig    `envPrefix:"MINIO_"`
		ProviderEndpoint         string         `env:"PROVIDER_ENDPOINT,required,notEmpty"`
		ProviderApiKey           string         `env:"PROVIDER_API_KEY,required,notEmpty"`
		ShutdownServerTimeoutSec int            `env:"SHUTDOWN_SERVER_TIMEOUT_SEC,required,notEmpty"`
	}

	DatabaseConfig struct {
		Driver   string `env:"DRIVER,required,notEmpty"`
		Host     string `env:"HOST,required,notEmpty"`
		Port     int    `env:"PORT,required,notEmpty"`
		Name     string `env:"NAME,required,notEmpty"`
		User     string `env:"USER,required,notEmpty"`
		Password string `env:"PASSWORD,required,notEmpty"`
	}

	MinioConfig struct {
		EndpointHost string `env:"ENDPOINT_HOST,required,notEmpty"`
		EndpointPort int    `env:"ENDPOINT_PORT,required,notEmpty"`

		BucketName string `env:"BUCKET_NAME,required,notEmpty"`
		User       string `env:"ROOT_USER,required,notEmpty"`
		Password   string `env:"ROOT_PASSWORD,required,notEmpty"`
		UseSSL     bool   `env:"USE_SSL,required,notEmpty"`
	}
)

func LoadConfig() (cfg Config, err error) {
	err = env.Parse(&cfg)
	return
}
