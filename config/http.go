package config

import (
	"fmt"
	"net"
	"os"
)

const (
	httpHost = "HTTP_HOST"
	httpPort = "HTTP_PORT"
)

type HTTPConfig interface {
	Address() string
}

type httpConfig struct {
	host string
	port string
}

// Address implements HTTPConfig.
func (cfg *httpConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}

func NewHTTPConfig() (HTTPConfig, error) {
	host := os.Getenv(httpHost)
	if len(host) == 0 {
		return nil, fmt.Errorf("env %v not found", httpHost)
	}

	port := os.Getenv(httpPort)
	if len(port) == 0 {
		return nil, fmt.Errorf("env %v not found", httpPort)
	}

	return &httpConfig{
		host: host,
		port: port,
	}, nil
}
