package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shresthashim/rest-api-golang/internal/config"
	"github.com/shresthashim/rest-api-golang/internal/http/handlers/task"
	"github.com/shresthashim/rest-api-golang/internal/storage/sqlite"
)

func main() {

	cfg := config.MustLoadConfig()

	sqliteStorage, err := sqlite.NewSQLiteStorage(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize SQLite storage: %v", err)
	}
	defer sqliteStorage.Db.Close()

	router := http.NewServeMux()

	router.HandleFunc("POST /api/tasks", task.New(sqliteStorage))

	readTimeout, err := time.ParseDuration(cfg.HTTP.ReadTimeout)
	if err != nil {
		log.Fatalf("Failed to parse read timeout: %v", err)
	}
	writeTimeout, err := time.ParseDuration(cfg.HTTP.WriteTimeout)
	if err != nil {
		log.Fatalf("Failed to parse write timeout: %v", err)
	}
	idleTimeout, err := time.ParseDuration(cfg.HTTP.IdleTimeout)
	if err != nil {
		log.Fatalf("Failed to parse idle timeout: %v", err)
	}
	shutdownTimeout, err := time.ParseDuration(cfg.HTTP.ShutdownTimeout)
	if err != nil {
		log.Fatalf("Failed to parse shutdown timeout: %v", err)
	}

	server := &http.Server{
		Addr:         cfg.HTTP.Addr,
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	slog.Info("Server started", slog.String("addr", cfg.HTTP.Addr))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err = server.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-done

	slog.Info("Server stopping...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)

	if err != nil {
		slog.Error("Server forced to shutdown:", slog.String("error", err.Error()))
	}

	slog.Info("Server exited properly")

}
