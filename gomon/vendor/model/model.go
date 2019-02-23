package model

import (
	"context"
	"controller/config"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type key string

const (
	hostKey     = key("hostKey")
	usernameKey = key("usernameKey")
	passwordKey = key("passwordKey")
	databaseKey = key("databaseKey")
)

var MongoDB *mongo.Database
var Ctx context.Context

func init() {
	Ctx = context.Background()
	Ctx, cancel := context.WithCancel(Ctx)
	defer cancel()

	MongoDB = nil

	Ctx = context.WithValue(Ctx, hostKey, config.Config.MongoSrv)
	Ctx = context.WithValue(Ctx, usernameKey, config.Config.MongoUsr)
	Ctx = context.WithValue(Ctx, passwordKey, config.Config.MongoPwd)
	Ctx = context.WithValue(Ctx, databaseKey, config.Config.MongoDb)

	uri := fmt.Sprintf(`mongodb://%s:%s@%s/%s?authSource=admin`,
		Ctx.Value(usernameKey),
		Ctx.Value(passwordKey),
		Ctx.Value(hostKey),
		Ctx.Value(databaseKey),
	)

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Error connecting Database: %v", err)
	}
	err = client.Connect(Ctx)
	if err != nil {
		log.Fatal("Error connecting Database: %v", err)
	}
	MongoDB = client.Database(config.Config.MongoDb)
}
