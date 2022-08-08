package main

import (
	"context"
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go-clickstream/internal/config"
	clients "go-clickstream/internal/pkg/clients/mongodb"
)

type DbEvent struct {
	DocumentKey   documentKey `bson:"documentKey"`
	OperationType string      `bson:"operationType"`
	FullDocument  interface{} `bson:"fullDocument"`
}
type documentKey struct {
	ID primitive.ObjectID `bson:"_id"`
}

var (
	configPath = flag.String("config-url", "./app.toml", "config file")
	userID     = flag.Int64("userID", 0, "userID for listening changes")
	eventID    = flag.Int64("eventID", 0, "eventID for listening changes")
)

func main() {
	flag.Parse()
	logger := logrus.New()
	logger.SetFormatter(&nested.Formatter{
		HideKeys:    true,
		FieldsOrder: []string{"component", "category"},
	})

	if *userID == 0 && *eventID == 0 {
		logger.Fatal("userID or eventID must be given")
	}
	if *userID < 0 && *eventID < 0 {
		logger.Fatal("userID or eventID must be > 0")
	}
	cfg := config.Config{}

	_, err := toml.DecodeFile(*configPath, &cfg)
	if err != nil {
		logger.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := clients.GetMongoConnect(ctx, cfg.MongoURL)
	if err != nil {
		logger.Fatal(err)
	}
	defer func() {
		if err = conn.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	err = conn.Ping(ctx, nil)
	if err != nil {
		logger.Fatal(err)
	}

	var waitGroup sync.WaitGroup

	db := conn.Database("clickstream")
	if db == nil {
		logger.Fatal("database clickstream not found")
	}

	condPipeline := getMatchConditions(*eventID, *userID)
	opts := options.ChangeStream().SetFullDocument(options.UpdateLookup)
	stream, err := db.Watch(context.TODO(), condPipeline, opts)
	if err != nil {
		panic(err)
	}

	waitGroup.Add(1)

	routineCtx, cancelFn := context.WithCancel(context.Background())
	_ = cancelFn

	go listenToDBChangeStream(logger, routineCtx, &waitGroup, stream)

	waitGroup.Wait()
}

func listenToDBChangeStream(logger *logrus.Logger,
	routineCtx context.Context,
	waitGroup *sync.WaitGroup,
	stream *mongo.ChangeStream,
) {
	defer stream.Close(routineCtx)
	defer waitGroup.Done()

	for stream.Next(routineCtx) {
		var event DbEvent
		if err := stream.Decode(&event); err != nil {
			panic(err)
		}
		if event.OperationType == "insert" {
			fmt.Println("Insert operation detected")
		} else if event.OperationType == "update" {
			fmt.Println("Update operation detected")
		} else if event.OperationType == "delete" {
			fmt.Println("Delete operation detected : Unable to pull changes as its record is deleted")
		}

		if event.OperationType == "insert" || event.OperationType == "update" {
			data, writeErr := bson.MarshalExtJSON(event.FullDocument, false, false)
			if writeErr != nil {
				logger.Fatal(writeErr)
			}
			logger.Info(string(data))

		}
	}
}

func getMatchConditions(eventID int64, userID int64) mongo.Pipeline {
	matchOperationType := bson.D{{"$match", bson.D{
		{
			Key: "$or", Value: []bson.D{
				{{"operationType", "insert"}},
				{{"operationType", "update"}},
				{{"operationType", "delete"}},
				{{"operationType", "replace"}},
			},
		},
	}}}

	res := mongo.Pipeline{matchOperationType}

	if userID != 0 {
		res = append(res, bson.D{{"$match", bson.D{
			{
				Key: "fullDocument.event.userID", Value: userID,
			},
		},
		}})
	}
	if eventID != 0 {
		res = append(res, bson.D{{"$match", bson.D{
			{
				Key: "fullDocument.id", Value: eventID,
			},
		},
		}})
	}
	return res
}