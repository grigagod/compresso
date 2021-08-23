package server

import (
	"context"
	"log"
	"net/http"

	_ "net/http/pprof"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/grigagod/compresso/internal/auth/config"
	"github.com/jmoiron/sqlx"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/grigagod/compresso/docs" // load API Docs files (Swagger)
)

type Server struct {
	cfg *config.Config
	db  *sqlx.DB
}

func NewServer(cfg *config.Config, db *sqlx.DB) *Server {
	return &Server{
		cfg: cfg,
		db:  db,
	}
}

// @title Auth service
// @version 1.0
// @description This is an auto-generated API Docs.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email podkidysh2002@gmail.com
func (s *Server) Run() {
	router := mux.NewRouter()

	s.MapHandlers(router)
	// enable pprof
	router.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)

	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://127.0.0.1:5000/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("#swagger-ui"),
	))

	srv := http.Server{
		Addr:         s.cfg.Server.Addr,
		WriteTimeout: s.cfg.Server.WriteTimeout,
		ReadTimeout:  s.cfg.Server.ReadTimeout,
		IdleTimeout:  s.cfg.Server.IdleTimeout,
		Handler:      router,
	}
	go func() {
		log.Printf("Server is listening on Addr: %s", s.cfg.Server.Addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal("Error starting Server: ", err)
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
		log.Println(err)
	}
	// Optionally, you could run srv.Shutdown in a goroutine and block on	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
}
