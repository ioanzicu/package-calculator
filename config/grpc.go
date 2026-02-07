package config

import (
	"errors"
	"net"
	"os"
)

const (
	grpcHost = "GRPC_HOST"
	grpcPort = "GRPC_PORT"
)

type GRPCCnfig interface {
	Address() string
}

type grpcConfig struct {
	host string
	port string
}

// Address implements GRPCCnfig.
func (cfg *grpcConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}

func NewGRPCConfig() (GRPCCnfig, error) {
	host := os.Getenv(grpcHost)
	if len(host) == 0 {
		return nil, errors.New("grpc host not found")
	}

	port := os.Getenv(grpcPort)
	if len(port) == 0 {
		return nil, errors.New("grpc port not found")
	}

	return &grpcConfig{
		host: host,
		port: port,
	}, nil
}
