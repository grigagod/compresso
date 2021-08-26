package main

import (
	"log"

	"github.com/grigagod/compresso/internal/auth/config"
	"github.com/grigagod/compresso/internal/auth/server"
	"github.com/grigagod/compresso/pkg/db/postgres"

	_ "github.com/grigagod/compresso/docs" // load API Docs files (Swagger)
)

// @title Auth service
// @version 1.0
// @description This is an auto-generated API Docs.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email podkidysh2002@gmail.com
func main() {
	log.Println("Starting auth server")

	authCfg, err := config.LoadConfig("./internal/auth/config/config-local")
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
	defer db.Close()

	s := server.NewAuthServer(authCfg, db)

	s.Run()
}
