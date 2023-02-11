package main

import (
	"flag"
	"fmt"
	"os"

	pkgEventService "github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/event/service"
	pkgHelpers "github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/helpers"
	pkgTaskRepository "github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/task/repository"
	pkgTaskService "github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/task/service"
	pkgPostgres "github.com/AssylzhanZharzhanov/axxonsoft-test-service/pkg/database/postgres"
	pkgRedis "github.com/AssylzhanZharzhanov/axxonsoft-test-service/pkg/database/redis"

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

	// Define our flags.
	//
	fs := flag.NewFlagSet("", flag.ExitOnError)
	_ = fs.String("http-addr", fmt.Sprintf(":%s", cfg.Port), "HTTP listen address")
	err = fs.Parse(os.Args[1:])
	if err != nil {
		logFatal(err)
	}

	// Setup database connection
	//
	db, err := pkgPostgres.NewConnection(cfg.DSN)
	if err != nil {
		logFatal(err)
	}

	redisClient, err := pkgRedis.NewRedisClient(cfg)
	if err != nil {
		logFatal(err)
	}

	//Repository layer.
	//
	taskRepository := pkgTaskRepository.NewRepository(db)
	taskRedisRepository := pkgTaskRepository.NewRedisRepository(redisClient)

	//Service layer.
	//
	eventService := pkgEventService.NewService()
	_ = pkgTaskService.NewService(eventService, taskRepository, taskRedisRepository)

}
