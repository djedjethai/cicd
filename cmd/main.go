package main

import (
	// "flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/djedjethai/generation0/pkg/config"
	"github.com/djedjethai/generation0/pkg/deleter"
	"github.com/djedjethai/generation0/pkg/getter"
	"github.com/djedjethai/generation0/pkg/handlers/rest"
	lgr "github.com/djedjethai/generation0/pkg/logger"
	"github.com/djedjethai/generation0/pkg/setter"
	storage "github.com/djedjethai/generation0/pkg/storage"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// var store = make(map[string]*storage.Node)
var encryptK = "PX9PHFrdn79ljrjLDZHlV1t+BdxHRFf5"
var postgresConfig = config.PostgresDBParams{
	Host:     "localhost",
	DbName:   "transactions",
	User:     "postgres",
	Password: "password",
}

// var port = ":8080"

// default value
var fileLoggerActive = false
var dbLoggerActive = false

// see jaeger: http://localhost:16686/
func main() {

	stdExporter, err := stdout.NewExporter(
		stdout.WithPrettyPrint(),
	)

	// jaegerEndpoint := "http://localhost:14268/api/traces"
	// within kubernetes for dev
	// jaegerEndpoint := "http://simplest-collector:14268/api/traces"
	// within kubernetes for prod
	jaegerEndpoint := "http://simple-prod-collector:14268/api/traces"
	serviceName := "golru"

	jaegerExporter, err := jaeger.NewRawExporter(
		jaeger.WithCollectorEndpoint(jaegerEndpoint),
		jaeger.WithProcess(jaeger.Process{
			ServiceName: serviceName,
		}))

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSyncer(stdExporter),
		sdktrace.WithSyncer(jaegerExporter),
	)

	otel.SetTracerProvider(tp)

	tr := otel.GetTracerProvider().Tracer("try")

	// TODO uncomment for k8s
	port := os.Getenv("PORT")

	// storage(infra layer)
	storagePtr := storage.NewStorage(2)

	// services(domain layer)
	setSrv := setter.NewSetter(storagePtr, tr)
	getSrv := getter.NewGetter(storagePtr, tr)
	delSrv := deleter.NewDeleter(storagePtr, tr)

	// in case the srv crash, when start back it will read the logger and recover its state
	logger, err := initializeTransactionLog(setSrv, delSrv, fileLoggerActive)
	if err != nil {
		log.Panic("FileLogger initialization failed")
	}
	defer logger.CloseFileLogger()

	dbLogger, err := initializeTransactionLogDB(setSrv, delSrv, dbLoggerActive)
	if err != nil {
		log.Panic("dbLogger initialization failed")
	}
	defer logger.CloseFileLogger()

	// handler(application layer)
	router := rest.Handler(setSrv, getSrv, delSrv, logger, dbLogger)

	fmt.Printf("***** Service listening on port %s *****", port)
	log.Fatal(http.ListenAndServe(port, router))
}

func initializeTransactionLogDB(setSrv setter.Setter, delSrv deleter.Deleter, active bool) (lgr.TransactionLogger, error) {
	var err error

	dbLogger, err := lgr.NewPostgresTransactionLogger(postgresConfig, active)
	if err != nil {
		return nil, fmt.Errorf("failed to create db event logger: %w", err)
	}

	if active {
		events, errors := dbLogger.ReadEvents()
		e, ok := lgr.Event{}, true

		for ok && err == nil {
			select {
			case err, ok = <-errors:
			case e, ok = <-events:
				switch e.EventType {
				case lgr.EventDelete:
					err = delSrv.Delete(e.Key)
				case lgr.EventPut:
					err = setSrv.Set(e.Key, []byte(e.Value))
				}

			}
		}

		dbLogger.Run()
	}

	return dbLogger, err

}

func initializeTransactionLog(setSrv setter.Setter, delSrv deleter.Deleter, active bool) (lgr.TransactionLogger, error) {
	var err error

	fileLogger, err := lgr.NewFileTransactionLogger("transaction.log", encryptK, active)
	if err != nil {
		return nil, fmt.Errorf("failed to create event logger: %w", err)
	}

	if active {
		events, errors := fileLogger.ReadEvents()
		e, ok := lgr.Event{}, true

		for ok && err == nil {
			select {
			case err, ok = <-errors: // retrieve any error
			case e, ok = <-events:
				switch e.EventType {
				case lgr.EventDelete:
					err = delSrv.Delete(e.Key)
				case lgr.EventPut:
					err = setSrv.Set(e.Key, []byte(e.Value))
				}
			}
		}

		fileLogger.Run()
	}

	return fileLogger, err
}
