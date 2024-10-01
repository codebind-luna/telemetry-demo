package mongorepository

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/codebind-luna/telemetry-demo/internal/interfaces"
	"github.com/codebind-luna/telemetry-demo/internal/repository/mongorepository/mongoutils"
	"github.com/codebind-luna/telemetry-demo/pkg/logger"
)

const (
	dataBase               = "expressions"
	calculationsCollection = "calculations"
)

var _ interfaces.Repository = (*mongorepo)(nil)

type mongorepo struct {
	client   *mongo.Client
	ctx      context.Context
	cancelFn context.CancelFunc
	logger   *logger.Logger
}

func (m *mongorepo) Close() {
	// Release resource
	mongoutils.Close(m.client, m.ctx, m.cancelFn, m.logger)
}

// Builds a mongo URI from environment variables
func getMongoURI() (string, error) {
	mongoURL, _ := os.LookupEnv("URL")

	if mongoURL != "" {
		return mongoURL, nil
	}

	mongoURL = "mongodb://"

	return mongoURL, nil
}

// New - create a new instance of the mongo repository
func New(
	logger *logger.Logger,
) (interfaces.Repository, error) {

	mongoURL, err := getMongoURI()
	if err != nil {
		return nil, err
	}

	// Get Client, Context, CancelFunc and
	// err from connect method.

	logger.Info("initializing connection to mongo cluster")
	client, ctx, cancel, err := mongoutils.Connect(mongoURL) //"mongodb://localhost:27017")
	if err != nil {
		return nil, err
	}

	// Ping mongoDB with Ping method
	logger.Info("initiating ping to confirm connection")
	err = mongoutils.Ping(client, ctx, logger)
	if err != nil {
		return nil, err
	}

	return &mongorepo{
		client:   client,
		ctx:      ctx,
		cancelFn: cancel,
		logger:   logger,
	}, nil
}

// calculationsCollection - get the calculations collection of expressions db
func (m *mongorepo) getcalculationsCollection() *mongo.Collection {
	return mongoutils.GetCollection(m.client, dataBase, calculationsCollection)
}
