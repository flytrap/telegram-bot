package app

import (
	"context"

	"github.com/flytrap/telegram-bot/internal/config"
	"github.com/redis/rueidis"
	"github.com/sirupsen/logrus"
)

type options struct {
	ConfigFile string
	Version    string
	UpdateDb   int64
	Index      string
}

type Option func(*options)

func SetConfigFile(s string) Option {
	return func(o *options) {
		o.ConfigFile = s
	}
}

func SetVersion(s string) Option {
	return func(o *options) {
		o.Version = s
	}
}

func SetUpdateDb(s int64) Option {
	return func(o *options) {
		o.UpdateDb = s
	}
}

func SetIndex(s string) Option {
	return func(o *options) {
		o.Index = s
	}
}

func InitStore() (rueidis.CoreClient, error) {
	client, err := rueidis.NewClient(rueidis.ClientOption{InitAddress: []string{config.C.Redis.Addr}, Password: config.C.Redis.Password})
	if err != nil {
		logrus.Warning(err)
		return nil, err
	}

	return client, nil
}

func RunIndex(ctx context.Context, opts ...Option) error {
	var o options
	for _, opt := range opts {
		opt(&o)
	}
	config.MustLoad(o.ConfigFile)
	config.PrintWithJSON()

	injector, err := BuildInjector()
	if err != nil {
		return err
	}

	if o.Index == "load" {
		injector.IndexManager.InitIndex(ctx)
		for _, lang := range config.C.Bot.Languages {
			injector.IndexManager.LoadData(ctx, lang)
		}
		return nil
	}
	if o.Index == "delete" {
		injector.IndexManager.DeleteAllIndex(ctx)
		return nil
	}

	if o.UpdateDb > 0 {
		injector.BotManager.UpdateGroupInfo(o.UpdateDb)
	} else {
		injector.IndexManager.InitIndex(ctx)
		go injector.BotManager.Start()
		return injector.GrpcServer.Run()
	}
	return nil
}
