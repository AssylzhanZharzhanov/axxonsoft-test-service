package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	pkgEventService "github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/event/service"
	pkgHelpers "github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/helpers"
	pkgPublisher "github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/publisher/service"
	pkgTaskRepository "github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/task/repository"
	pkgTaskService "github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/task/service"
	pkgTaskEndpoints "github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/task/transport"
	pkgPostgres "github.com/AssylzhanZharzhanov/axxonsoft-test-service/pkg/database/postgres"
	pkgRedis "github.com/AssylzhanZharzhanov/axxonsoft-test-service/pkg/database/redis"
	pkgRabbitMQ "github.com/AssylzhanZharzhanov/axxonsoft-test-service/pkg/rabbitmq"

	kitzapadapter "github.com/go-kit/kit/log/zap"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
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
	httpAddr := fs.String("http-addr", fmt.Sprintf(":%s", cfg.Port), "HTTP listen address")
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

	_, err = pkgRabbitMQ.DeclareBinding(amqpChan, cfg.ExchangeName, cfg.ExchangeKind)
	if err != nil {
		logFatal(err)
	}

	// Repository layer.
	//
	taskRepository := pkgTaskRepository.NewRepository(db)
	taskRedisRepository := pkgTaskRepository.NewRedisRepository(redisClient)

	// Service layer.
	//
	publisher := pkgPublisher.NewService(amqpConn, amqpChan, cfg.ExchangeName, cfg.ExchangeKind, logger)
	eventService := pkgEventService.NewService(publisher, logger)
	taskService := pkgTaskService.NewService(eventService, taskRepository, taskRedisRepository, logger)

	// Endpoints layer.
	//
	taskEndpoints := pkgTaskEndpoints.NewEndpoints(taskService, logger)

	r := mux.NewRouter()
	{
		pkgTaskEndpoints.RegisterRoutersV1(r, taskEndpoints, logger)
	}
	http.Handle("/", accessControl(r, "*"))

	// This function just sits and waits for ctrl-C.
	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()
	go func() {
		_ = logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, nil)
	}()
	_ = logger.Log("exit", <-errs)
}

func accessControl(h http.Handler, allowedOrigins string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedHeaders := "Accept, Origin, Content-Type, Content-Length, Authorization, Firebase-Authorization, X-Request-Server-TimeZone, X-Time-Zone, x-requested-with, Idempotency-Key"

		w.Header().Set("Access-Control-Allow-Origin", allowedOrigins)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
