package worker

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"

	"github.com/codebind-luna/telemetry-demo/internal/interfaces"
	"github.com/codebind-luna/telemetry-demo/pkg/logger"
	"github.com/codebind-luna/telemetry-demo/pkg/publisher/pubsubconstants"
	wotel "github.com/voi-oss/watermill-opentelemetry/pkg/opentelemetry"
)

type Worker interface {
	Start()
	Stop(ctx context.Context)
	Healthy(ctx context.Context) bool
}

var _ Worker = (*workerImpl)(nil)

type workerImpl struct {
	svc    interfaces.WorkerService
	router *message.Router
	log    *logger.Logger
}

// Healthy - check if the worker is healthy
func (w *workerImpl) Healthy(ctx context.Context) bool {
	return w.router.IsRunning()
}

// Start - start the worker
func (w *workerImpl) Start() {
	go func() {
		w.router.Run(context.Background())
	}()
}

// Stop - stop the worker
func (w *workerImpl) Stop(ctx context.Context) {
	w.router.Close()
}

// setup - configure the worker
func (w *workerImpl) setup(
	expressionSubscriber message.Subscriber,
) {
	// Subscribe to the expression-handler topic
	w.router.AddNoPublisherHandler(
		"expression-handler",
		pubsubconstants.ExpressionTopicName,
		expressionSubscriber,
		w.handleExpression,
	)
}

// New - create a new instance
func New(
	loggerInstance *logger.Logger,
	svc interfaces.WorkerService,
	eventsProcessor interfaces.EventProcessor,
) Worker {
	loggerInstance.Info("initializing worker")

	// Initialize amqp router
	router, err := message.NewRouter(
		message.RouterConfig{},
		nil,
	)
	if err != nil {
		loggerInstance.Fatal("problem initializing router")
	}

	// Add middleware
	router.AddMiddleware(
		middleware.Recoverer,
	)

	router.AddMiddleware(wotel.Trace())

	w := &workerImpl{
		log:    loggerInstance,
		svc:    svc,
		router: router,
	}

	w.setup(
		eventsProcessor.GetExpressionSubscriber(),
	)

	return w
}
