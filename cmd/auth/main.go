package main

import (
	"log"

	"github.com/grigagod/compresso/internal/auth/config"
	"github.com/grigagod/compresso/internal/auth/server"
	"github.com/grigagod/compresso/pkg/db/postgres"
)

func main() {
	log.Println("Starting auth server")

	cfg, err := config.LoadConfig("./internal/auth/config/config-local")
	if err != nil {
		log.Fatal(err)
	}

	pgCfg, err := postgres.LoadConfig("./config/config-postgres")
	if err != nil {
		log.Fatal(err)
	}

	db, err := postgres.NewPsqlDB(pgCfg)
	if err != nil {
		log.Fatal("Postgres connection failed:", err)
	}

	server := server.NewServer(cfg, db)

	server.Run()
}
