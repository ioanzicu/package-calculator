package config

import (
	"fmt"
	"os"
	"time"
)

const gracefulShutdownTimeoutSec = "GRACEFUL_SHUTDOWN_TIMEOUT_SECONDS"

type GracefulShutdownConfig interface {
	Timeout() time.Duration
}

type gracefulShutdownConfig struct {
	duration time.Duration
}

// TimeoutSeconds implements GracefulShutdownConfig.
func (g *gracefulShutdownConfig) Timeout() time.Duration {
	return g.duration
}

func NewGracefulShutdownConfig() (GracefulShutdownConfig, error) {
	durationStr := os.Getenv(gracefulShutdownTimeoutSec)
	if len(durationStr) == 0 {
		return nil, fmt.Errorf("env %v not found", gracefulShutdownTimeoutSec)
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return nil, err
	}

	return &gracefulShutdownConfig{
		duration: duration,
	}, nil
}
