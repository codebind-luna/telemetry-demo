package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/codebind-luna/telemetry-demo/internal/backend"
	"github.com/codebind-luna/telemetry-demo/internal/eventprocessor"
	"github.com/codebind-luna/telemetry-demo/internal/interfaces"
	"github.com/codebind-luna/telemetry-demo/internal/repository"
	workerService "github.com/codebind-luna/telemetry-demo/internal/services/worker-service"
	"github.com/codebind-luna/telemetry-demo/pkg/logger"
	"github.com/codebind-luna/telemetry-demo/pkg/tracermetrics"
	"github.com/codebind-luna/telemetry-demo/pkg/worker"
)

const (
	DefaultShutdownContextWaitSeconds = 5
)

func main() {
	logger := logger.NewLogger()
	repoType, repoErr := interfaces.ParseRepository("mongo")
	if repoErr != nil {
		logger.Fatal(repoErr.Error())
	}
	repo, rErr := repository.New(repoType, logger)
	if rErr != nil {
		logger.Fatal(rErr.Error())
	}

	eventBusType, busErr := interfaces.ParseMessageBusType("amqp")
	if busErr != nil {
		logger.Fatal(busErr.Error())
	}

	backend := backend.NewBackend(logger, repo)

	workerSvc := workerService.NewService(logger, backend)

	eventProcessor, createEventProcessorErr := eventprocessor.New(eventBusType)
	if createEventProcessorErr != nil {
		logger.Fatal(createEventProcessorErr.Error())
	}

	worker := worker.New(logger, workerSvc, eventProcessor)

	exporterType, createErr := interfaces.ParseTracerExporterType("jaeger")
	if createErr != nil {
		logger.Fatal(createErr.Error())
	}

	collector := tracermetrics.New("expression-worker", "1.0.0", logger, true, exporterType)

	worker.Start()
	collector.Start()

	// Create shutdown chan
	startShutdown := make(chan struct{})
	// Gracefully shut down
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigs
		signal.Stop(sigs)
		log.Printf("received signal %v\n", sig)
		close(startShutdown)

		// Allow pressing Ctrl+C again to exit, otherwise the developer must manually kill the process
		if sig == syscall.SIGINT {
			sigs2 := make(chan os.Signal, 1)
			signal.Notify(sigs2, syscall.SIGINT)
			log.Println("press Ctrl+C again to exit")
			<-sigs2
			os.Exit(0)
		}
	}()

	// Block until a signal is received
	<-startShutdown
	// Create the cancelation context
	ctx, cancel := context.WithTimeout(
		context.Background(),
		DefaultShutdownContextWaitSeconds*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	// Stop worker
	worker.Stop(ctx)
	// Stop trace metrics collector
	collector.Stop(ctx)
}
