package config

import (
	"errors"
	"os"
)

const dsnEnv = "PG_DSN"

type PGConfig interface {
	DSN() string
}

type pgConfig struct {
	dsn string
}

// DSN implements PGConfig.
func (cfg *pgConfig) DSN() string {
	return cfg.dsn
}

func NewPGConfig() (PGConfig, error) {
	dsn := os.Getenv(dsnEnv)
	if len(dsn) == 0 {
		return nil, errors.New("pg dsn not found")
	}

	return &pgConfig{
		dsn: dsn,
	}, nil
}
