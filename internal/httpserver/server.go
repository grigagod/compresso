package httpserver

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/grigagod/compresso/pkg/logger"
)

// ListenAndServe starts http server and gracefully stops it
func ListenAndServe(cfg *Config, router *chi.Mux, log logger.Logger) {
	if cfg.Pprof {
		MapPprofHandler(router)
	}

	if cfg.Swagger {
		MapSwaggerHandler(router, cfg.SwaggerUrl)
	}

	srv := http.Server{
		Addr:         cfg.Addr,
		WriteTimeout: cfg.WriteTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		Handler:      router,
	}

	// channel for graceful shutdown
	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Run the server
	go func() {
		log.Infof("Server is listening on Addr: %s", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting Server: ", err)
		}
	}()

	// Listen for syscall signals
	<-done
	log.Infof("Server is shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed: %+v", err)

	}

	log.Infof("Server Exited Properly")
}
