package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/shresthashim/rest-api-golang/internal/config"
)

func main() {

	cfg := config.MustLoadConfig()

	var err error

	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to my API!"))
	})

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

	server := &http.Server{
		Addr:         cfg.HTTP.Addr,
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	fmt.Println("Server started on", cfg.HTTP.Addr)

	err = server.ListenAndServe()

	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
