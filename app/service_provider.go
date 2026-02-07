package app

import (
	"context"
	"ignis/config"
	"ignis/internal/adapter/db"
	"ignis/internal/domain"
	"ignis/internal/service"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type serviceProvider struct {
	httpConfig             config.HTTPConfig
	gracefulShutdownConfig config.GracefulShutdownConfig
	pgConfig               config.PGConfig
	pgPool                 *pgxpool.Pool
	dbRepository           db.Repository
	packageCalculator      domain.PackageCalculator
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := config.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get http config: %s", err.Error())
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

func (s *serviceProvider) GracefulShutdownConfig() config.GracefulShutdownConfig {
	if s.gracefulShutdownConfig == nil {
		cfg, err := config.NewGracefulShutdownConfig()
		if err != nil {
			log.Fatalf("failed to get gracefulShutdown config: %s", err.Error())
		}

		s.gracefulShutdownConfig = cfg
	}

	return s.gracefulShutdownConfig
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) PGPool(ctx context.Context) *pgxpool.Pool {
	if s.pgPool == nil {
		pool, err := pgxpool.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to connect to database: %s", err.Error())
		}

		err = pool.Ping(ctx)
		if err != nil {
			log.Fatalf("failed to ping database: %s", err.Error())
		}

		s.pgPool = pool
	}

	return s.pgPool
}

func (s *serviceProvider) DBRepository(ctx context.Context) db.Repository {
	if s.dbRepository == nil {
		s.dbRepository = db.NewRepository(s.PGPool(ctx))
	}

	return s.dbRepository
}

func (s *serviceProvider) PackageCalculator() domain.PackageCalculator {
	if s.packageCalculator == nil {
		s.packageCalculator = service.NewPackageCalculatorService()
	}

	return s.packageCalculator
}
