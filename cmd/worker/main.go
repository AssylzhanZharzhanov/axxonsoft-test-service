package main

import (
	"context"
	pkgConsumer "github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/consumer/service"
	pkgHelpers "github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/helpers"
	pkgTaskRepository "github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/task/repository"
	pkgTaskRunnerService "github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/taskrunner/service"
	pkgPostgres "github.com/AssylzhanZharzhanov/axxonsoft-test-service/pkg/database/postgres"
	pkgRabbitMQ "github.com/AssylzhanZharzhanov/axxonsoft-test-service/pkg/rabbitmq"
	"os"

	kitzapadapter "github.com/go-kit/kit/log/zap"
	"github.com/go-kit/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {

	// Create a single logger, which we'll use and give to other components.
	//
	zapLogger, _ := zap.NewProduction()
	defer func() {
		_ = zapLogger.Sync()
	}()

	var logger log.Logger
	logger = kitzapadapter.NewZapSugarLogger(zapLogger, zapcore.InfoLevel)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)
	// Logging helper function
	logFatal := func(err error) {
		_ = logger.Log("err", err)
		os.Exit(1)
	}

	cfg, err := pkgHelpers.LoadConfig()
	if err != nil {
		logFatal(err)
	}

	// Setup database connection
	//
	db, err := pkgPostgres.NewConnection(cfg.DSN)
	if err != nil {
		logFatal(err)
	}

	// Setup amqp connection
	amqpConn, err := pkgRabbitMQ.NewRabbitMQConnection(cfg.RabbitMQURI)
	if err != nil {
		logFatal(err)
	}
	defer amqpConn.Close()

	amqpChan, err := amqpConn.Channel()
	if err != nil {
		logFatal(err)
	}
	defer amqpChan.Close()

	_, err = pkgRabbitMQ.DeclareQueue(amqpChan, cfg.QueueName)
	if err != nil {
		logFatal(err)
	}

	// Repository layer
	//
	taskRepository := pkgTaskRepository.NewRepository(db)

	// Service layer
	//
	taskRunnerService := pkgTaskRunnerService.NewTaskRunnerService(taskRepository, logger)
	consumerService := pkgConsumer.NewConsumerService(taskRunnerService, amqpChan, cfg.QueueName, cfg.Consumer, logger)

	var (
		ctx = context.Background()
	)
	err = consumerService.Consume(ctx)
	if err != nil {
		logFatal(err)
	}
}
