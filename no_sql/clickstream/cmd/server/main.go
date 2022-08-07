package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/labstack/echo/v4"

	"go-clickstream/internal/api/handler/get_event"
	"go-clickstream/internal/api/handler/send_event"
	"go-clickstream/internal/api/handler/upload_event"
	"go-clickstream/internal/pkg/clients/mongodb"
	events2 "go-clickstream/internal/repository/events"
	"go-clickstream/internal/usecase/events"

	"go-clickstream/internal/config"
)

var (
	configPath = flag.String("config-url", "./app.toml", "config file")
)

func main() {
	flag.Parse()

	cfg := config.Config{}

	_, err := toml.DecodeFile(*configPath, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := clients.GetMongoConnect(ctx, cfg.MongoURL)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = conn.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	err = conn.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	db := conn.Database("clickstream")
	if db == nil {
		log.Fatal("database clickstream not found")
	}

	eventsRepo := events2.NewRepository(conn, db, "events")
	eventsUsecase := events.NewEventsUsecase(eventsRepo)
	sendEventHandler := send_event.NewHandler(eventsUsecase)
	getEventHandler := get_event.NewHandler(eventsUsecase)
	uploadEventHandler := upload_event.NewHandler(eventsUsecase)

	e := echo.New()

	e.POST("/event/send", sendEventHandler.Handle)
	e.GET("/event/:eventID", getEventHandler.Handle)
	e.POST("/event/upload", uploadEventHandler.Handle)

	err = e.Start(cfg.BindAddr)
	if err != nil {
		log.Fatal(err)
	}
}
