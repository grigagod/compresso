package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/grigagod/compresso/internal/config"
	"github.com/grigagod/compresso/internal/storage"
	videoCfg "github.com/grigagod/compresso/internal/video/config"
	"github.com/grigagod/compresso/internal/video/convsvc"
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

	var cfg videoCfg.Converter
	err = config.LoadConfig(&cfg, config.GetConfigPath("videosvc", os.Getenv("config-videosvc")))
	if err != nil {
		logger.Fatal(err)
	}

	db, err := postgres.NewPsqlDB(cfg.DB.Host, cfg.DB.Port, cfg.DB.User,
		cfg.DB.DBName, cfg.DB.Password, cfg.DB.Driver)
	if err != nil {
		logger.Fatalf("Postgres connection failed:", err)
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
	defer ch.Close()

	svc := convsvc.NewService(ch, db, storage, logger)

	ctx, cancel := context.WithCancel(context.Background())

	// Gracefully shutdown
	go func(cancel context.CancelFunc) {
		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		<-done
		logger.Infof("Shutting down gracefully")

		cancel()
	}(cancel)

	if err := svc.Run(ctx, &cfg.QueueReadConfig); err != nil {
		logger.Fatal("Service failed to run:", err)
	}
}