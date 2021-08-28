package main

import (
	"log"
	"os"

	"github.com/grigagod/compresso/internal/auth/config"
	"github.com/grigagod/compresso/internal/auth/server"
	"github.com/grigagod/compresso/pkg/db/postgres"
	"github.com/grigagod/compresso/pkg/utils"

	_ "github.com/grigagod/compresso/docs/auth" // load API Docs files (Swagger)
)

// @title Auth service
// @version 1.0
// @description This is an auto-generated API Docs.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email podkidysh2002@gmail.com
func main() {
	log.Println("Starting auth server")

	authCfg, err := config.LoadConfig(utils.GetConfigPath("auth", os.Getenv("config-auth")))
	if err != nil {
		log.Fatal(err)
	}

	db, err := postgres.NewPsqlDB(authCfg.DB.Host, authCfg.DB.Port, authCfg.DB.User,
		authCfg.DB.DbName, authCfg.DB.Password, authCfg.DB.Driver)
	if err != nil {
		log.Fatal("Postgres connection failed:", err)
	}
	defer db.Close()

	s := server.NewAuthServer(authCfg, db)

	s.Run()
}
