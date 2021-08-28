package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	Cfg    Config
	Router *mux.Router
	Db     *sqlx.DB
	Logger *log.Logger
}

func (s *Server) Start() {
	s.MapDefaultHandlers()

	srv := http.Server{
		Addr:         s.Cfg.Addr,
		WriteTimeout: s.Cfg.WriteTimeout,
		ReadTimeout:  s.Cfg.ReadTimeout,
		IdleTimeout:  s.Cfg.IdleTimeout,
		Handler:      s.Router,
	}
	// channel for graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// gracefully shutdown server on SIGINT, SIGTERM
	go func() {

		s.Logger.Printf("Server is listening on Addr: %s", s.Cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.Logger.Fatal("Error starting Server: ", err)
		}
	}()

	<-done
	s.Logger.Println("Server is shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		s.Logger.Fatalf("Server Shutdown Failed:%+v", err)

	}

	s.Logger.Print("Server Exited Properly")

}
