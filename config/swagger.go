package config

import (
	"errors"
	"net"
	"os"
)

const (
	swaggerHost = "SWAGGER_HOST"
	swaggerPort = "SWAGGER_PORT"
)

type SwaggerConfig interface {
	Address() string
}

type swaggerConfig struct {
	host string
	port string
}

// Address implements SwaggerConfig.
func (cfg *swaggerConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}

func NewSwaggerConfig() (SwaggerConfig, error) {
	host := os.Getenv(swaggerHost)
	if len(host) == 0 {
		return nil, errors.New("swagger host not found")
	}

	port := os.Getenv(swaggerPort)
	if len(port) == 0 {
		return nil, errors.New("swagger port not found")
	}

	return &swaggerConfig{
		host: host,
		port: port,
	}, nil
}
