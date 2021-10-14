package main

import (
	"log"
	"os"

	"github.com/grigagod/compresso/internal/config"
	"github.com/grigagod/compresso/internal/storage"
	videoCfg "github.com/grigagod/compresso/internal/video/config"
	"github.com/grigagod/compresso/internal/video/server"
	"github.com/grigagod/compresso/pkg/db/aws"
	"github.com/grigagod/compresso/pkg/db/postgres"
	"github.com/grigagod/compresso/pkg/logger"
	"github.com/grigagod/compresso/pkg/rmq"
)

func main() {
	logger, err := logger.NewWrappedLogger(logger.GetLoggerConfig())
	if err != nil {
		log.Fatal(err)
	}

	var cfg videoCfg.API
	err = config.LoadConfig(&cfg, config.GetConfigPath("videoapi", os.Getenv("config-videoapi")))
	if err != nil {
		logger.Fatal(err)
	}

	db, err := postgres.NewPsqlDB(cfg.DB.Host, cfg.DB.Port, cfg.DB.User,
		cfg.DB.DbName, cfg.DB.Password, cfg.DB.Driver)
	if err != nil {
		logger.Fatal("Postgres connection failed:", err)
	}
	defer db.Close()

	s3client, err := aws.NewClientWithSharedCredentials("./.aws/credentials", "test")
	if err != nil {
		logger.Fatal("AWS S3 session failed:", err)
	}
	storage := storage.NewAWSStorage(cfg.Storage, s3client)
	logger.Infof("AWS Bucket: %s", cfg.Storage.Bucket)

	ch, err := rmq.NewChannelFromConfig(&cfg.RMQ)
	if err != nil {
		logger.Fatal("RMQ connection failed:", err)
	}
	pub := rmq.NewPublisher(ch)

	s := server.NewVideoServer(&cfg.APIsvc, db, storage, pub, logger)
	s.MapHandlers()
	s.ListenAndServe(&cfg.HTTP)
}
