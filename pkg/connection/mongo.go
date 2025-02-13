package connection

import (
	"SangXanh/pkg/config"
	"SangXanh/pkg/log"
	"context"
	"github.com/samber/do/v2"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

func NewMongoDatabase(di do.Injector) (*mongo.Database, error) {
	ctx := context.Background()
	conf := do.MustInvoke[config.Mongo](di)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.URI).SetMonitor(&event.CommandMonitor{
		Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
			log.Info(startedEvent.Command)
		},
		Succeeded: nil,
		Failed:    nil,
	}))
	if err != nil {
		log.Errorw("failed to connect to mongo", "error", err)
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}
	return client.Database(conf.Database), nil
}
