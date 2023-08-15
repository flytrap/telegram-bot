package app

import (
	"context"

	"github.com/flytrap/telegram-bot/internal/config"
)

type options struct {
	ConfigFile string
	Version    string
	UpdateDb   bool
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

func SetUpdateDb(s bool) Option {
	return func(o *options) {
		o.UpdateDb = s
	}
}

// func InitStore() (*redis.Store, error) {
// 	cfg := config.C.Redis
// 	c := redis.Config{}
// 	copier.Copy(&c, cfg)
// 	store := redis.NewStore(&c)
// 	return store, nil
// }

func Run(ctx context.Context, opts ...Option) error {
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

	if o.UpdateDb {
		injector.BotManager.UpdateGroupInfo()
	} else {
		injector.BotManager.Start()
	}
	return nil
}
