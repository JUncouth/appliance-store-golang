// @Version 1.0.0
// @Title StoreAPI
// @Description REST API for an appliance store.
// @Server https://xoxotyh-appliance-store appliance store
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpapi "storeapi/src/adapters/http"
	"storeapi/src/adapters/postgres"
	"storeapi/src/internal/application"
	"storeapi/src/internal/config"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	store, err := postgres.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("database startup failed", "error", err)
		os.Exit(1)
	}
	defer store.Close()
	clients := application.NewClientService(postgres.NewClientRepository(store))
	products := application.NewProductService(postgres.NewProductRepository(store))
	suppliers := application.NewSupplierService(postgres.NewSupplierRepository(store))
	images := application.NewImageService(postgres.NewImageRepository(store))
	handlers := httpapi.NewHandler(clients, products, suppliers, images)
	server := &http.Server{Addr: ":" + cfg.HTTPPort, Handler: httpapi.NewRouter(handlers), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 30 * time.Second, IdleTimeout: 60 * time.Second}
	go func() {
		logger.Info("StoreAPI started", "port", cfg.HTTPPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("http server failed", "error", err)
			cancel()
		}
	}()
	<-ctx.Done()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
	}
	logger.Info("StoreAPI stopped")
}
