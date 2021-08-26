package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	Cfg    *Config
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

	go func() {

		s.Logger.Printf("Server is listening on Addr: %s", s.Cfg.Addr)
		if err := srv.ListenAndServe(); err != nil {
			s.Logger.Fatal("Error starting Server: ", err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.

	err := srv.Shutdown(ctx)
	if err != nil {
		s.Logger.Println(err)
	}
	// Optionally, you could run srv.Shutdown in a goroutine and block on	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	s.Logger.Println("shutting down")

}
