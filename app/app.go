package app

import (
	"context"
	"ignis/closer"
	"ignis/config"
	"ignis/internal/adapter/api"
	"log"
	"net/http"
	"sync"
	"time"
)

type App struct {
	serviceProvider *serviceProvider
	httpServer      *http.Server
}

func NewApp() (*App, error) {
	a := &App{}
	err := a.initConfig()
	if err != nil {
		return nil, err
	}

	a.initServiceProvider()

	err = a.initHTTPServer()
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	defer func() {
		closer.Wait()
	}()

	wg := sync.WaitGroup{}

	wg.Go(func() {
		err := a.runHTTPServer()
		if err != nil {
			log.Printf("failed to run HTTP server: %v\n", err)
		}
	})

	wg.Wait()

	return nil
}

func (a *App) initConfig() error {
	err := config.Load(".env")
	if err != nil {
		return nil
	}

	return nil
}

func (a *App) initServiceProvider() {
	a.serviceProvider = newServiceProvider()
}

func (a *App) initHTTPServer() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", api.RootHandler)

	// Calculator handler
	calcRepo := a.serviceProvider.DBRepository(context.Background())
	calculatorHandler := api.NewCalculatorHandler(a.serviceProvider.PackageCalculator(), calcRepo)
	mux.HandleFunc("/api/v1/calculate", calculatorHandler.Calculate)
	mux.HandleFunc("/api/v1/history", calculatorHandler.History)

	mux.HandleFunc("/healthz", api.HealthHandler)

	a.httpServer = &http.Server{
		Addr:         a.serviceProvider.HTTPConfig().Address(),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	closer.Add(func(ctx context.Context) error {
		return a.httpServer.Shutdown(ctx)
	})

	return nil
}

func (a *App) runHTTPServer() error {
	log.Printf("HTTP Server is running on %v\n", a.serviceProvider.HTTPConfig().Address())

	err := a.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
