package transport

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/codebind-luna/telemetry-demo/internal/handlers"
	"github.com/codebind-luna/telemetry-demo/internal/interfaces"
	"github.com/codebind-luna/telemetry-demo/internal/middleware"
	"github.com/codebind-luna/telemetry-demo/pkg/constants"
	"github.com/codebind-luna/telemetry-demo/pkg/logger"
)

type Server interface {
	Start()
	Stop(ctx context.Context)
}

var _ Server = (*serverImp)(nil)

type serverImp struct {
	port   int
	logger *logger.Logger
	srv    http.Server
}

func NewServer(logger *logger.Logger, port int, srvc interfaces.ExpressionsService) Server {
	expressionsHandler := handlers.NewExpressionsHandler(logger, srvc)

	// Create a new ServeMux instance
	mux := http.NewServeMux()

	// Register routes with the ServeMux
	mux.HandleFunc("/calculate", middleware.OTelMiddleware(constants.ServiceName, expressionsHandler.Calculate))

	// Define a route with a path parameter
	mux.HandleFunc("/expressions/{id}", middleware.OTelMiddleware(constants.ServiceName, expressionsHandler.Get))

	return &serverImp{
		logger: logger,
		port:   port,
		srv: http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		},
	}
}

func (s *serverImp) Start() {
	s.logger.Infof("starting internal http server on :%d", s.port)

	go func() {
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Fatalf("error starting http server on port %d : %s", s.port, err.Error())
		}
	}()
}

func (s *serverImp) Stop(ctx context.Context) {
	s.logger.Infof("initializing stop of http server on port %d", s.port)

	if err := s.srv.Shutdown(ctx); err != nil {
		s.logger.Errorf("error stopping http server: %s", err.Error())
	}
}
