package mongoutils

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/codebind-luna/telemetry-demo/pkg/logger"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	// "go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

var (
	tUUID       = reflect.TypeOf(uuid.UUID{})
	uuidSubtype = byte(0x04)
)

// uuidEncodeValue - encoder for UUID to a bson value
func uuidEncodeValue(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.Type() != tUUID {
		return bsoncodec.ValueEncoderError{Name: "uuidEncodeValue", Types: []reflect.Type{tUUID}, Received: val}
	}
	b := val.Interface().(uuid.UUID)
	return vw.WriteBinaryWithSubtype(b[:], uuidSubtype)
}

// uuidDecodeValue - decode a bson value to a UUID
func uuidDecodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() || val.Type() != tUUID {
		return bsoncodec.ValueDecoderError{Name: "uuidDecodeValue", Types: []reflect.Type{tUUID}, Received: val}
	}

	var data []byte
	var subtype byte
	var err error
	switch vrType := vr.Type(); vrType {
	case bson.TypeBinary:
		data, subtype, err = vr.ReadBinary()
		if subtype != uuidSubtype {
			return fmt.Errorf("unsupported binary subtype %v for UUID", subtype)
		}
	case bson.TypeNull:
		err = vr.ReadNull()
	case bson.TypeUndefined:
		err = vr.ReadUndefined()
	default:
		return fmt.Errorf("cannot decode %v into a UUID", vrType)
	}

	if err != nil {
		return err
	}
	uuid2, err := uuid.FromBytes(data)
	if err != nil {
		return err
	}
	val.Set(reflect.ValueOf(uuid2))
	return nil
}

// Close - this is a user defined method to close resources.
// This method closes mongoDB connection and cancel context.
func Close(
	client *mongo.Client,
	ctx context.Context,
	cancel context.CancelFunc,
	logger *logger.Logger,
) {

	// CancelFunc to cancel to context
	defer cancel()

	// client provides a method to close
	// a mongoDB connection.
	defer func() {

		// client.Disconnect method also has deadline.
		// returns error if any,
		if err := client.Disconnect(ctx); err != nil {
			logger.Fatal(err)
		}
	}()
}

// Connect - this is a user defined method that returns mongo.Client,
// context.Context, context.CancelFunc and error.
// mongo.Client will be used for further database operation.
// context.Context will be used set deadlines for process.
// context.CancelFunc will be used to cancel context and
// resource associated with it.
func Connect(
	uri string,
) (
	*mongo.Client,
	context.Context,
	context.CancelFunc,
	error,
) {

	// ctx will be used to set deadline for process, here
	// deadline will of 30 seconds.
	ctx, cancel := context.WithTimeout(context.Background(),
		30*time.Second)

	mongoRegistry := bson.NewRegistry()
	mongoRegistry.RegisterTypeEncoder(tUUID, bsoncodec.ValueEncoderFunc(uuidEncodeValue))
	mongoRegistry.RegisterTypeDecoder(tUUID, bsoncodec.ValueDecoderFunc(uuidDecodeValue))

	// mongo.Connect return mongo.Client method
	opts := options.Client()
	// opts.Monitor = otelmongo.NewMonitor()
	opts.ApplyURI(uri)
	client, err := mongo.Connect(ctx, opts)
	return client, ctx, cancel, err
}

// Ping - this is a user defined method that accepts
// mongo.Client and context.Context
// This method used to ping the mongoDB, return error if any.
func Ping(client *mongo.Client, ctx context.Context, logger *logger.Logger) error {

	// mongo.Client has Ping to ping mongoDB, deadline of
	// the Ping method will be determined by cxt
	// Ping method return error if any occurred, then
	// the error can be handled.
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	logger.Info("pinged successfully")
	return nil
}

// GetCollection - gets a collection based on the name that has been passed in
func GetCollection(client *mongo.Client, database, collection string) *mongo.Collection {
	return client.Database(database).Collection(collection)
}
