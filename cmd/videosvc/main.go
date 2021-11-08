package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/grigagod/compresso/internal/storage"
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

	dbCfg, err := postgres.GetConfigFromEnv()
	if err != nil {
		logger.Fatal(err)
	}

	storageCfg, err := storage.GetConfigFromEnv()
	if err != nil {
		logger.Fatal(err)
	}

	rmqCfg, err := rmq.GetConfigFromEnv()
	if err != nil {
		logger.Fatal(err)
	}

	qrCfg, err := rmq.GetQueueReadConfigFromEnv()
	if err != nil {
		logger.Fatal(err)
	}

	db, err := postgres.NewPsqlDB(dbCfg.URL, dbCfg.Driver, dbCfg.MaxOpenConns,
		dbCfg.MaxIdleConns, dbCfg.ConnMaxLifetime, dbCfg.ConnMaxIdleTime)
	if err != nil {
		logger.Fatal("Postgres connection failed:", err)
	}
	defer db.Close()

	s3client, err := aws.NewClientWithEnvCredentials()
	if err != nil {
		logger.Fatal("AWS S3 session failed:", err)
	}

	storage := storage.NewAWSStorage(storageCfg, s3client)

	ch, err := rmq.NewChannelFromConfig(rmqCfg)
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
		logger.Infof("Received cancel signal, start shutdown")

		cancel()
	}(cancel)

	if err := svc.Run(ctx, qrCfg); err != nil {
		logger.Fatal("Service failed to run:", err)
	}
}
