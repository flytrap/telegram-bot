package app

import (
	"context"
	"encoding/json"

	"github.com/flytrap/telegram-bot/internal/config"
	"github.com/flytrap/telegram-bot/pkg/redis"
	"github.com/jinzhu/copier"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/redis/rueidis"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
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

func InitIndex() (rueidis.CoreClient, error) {
	client, err := rueidis.NewClient(rueidis.ClientOption{InitAddress: []string{config.C.Redis.Addr}, Password: config.C.Redis.Password, SelectDB: config.C.Redis.DB})
	if err != nil {
		logrus.Warning(err)
		return nil, err
	}

	return client, nil
}

func InitStore() (*redis.Store, error) {
	c := redis.Config{}
	copier.Copy(&c, config.C.Redis)
	store := redis.NewStore(&c)
	return store, nil
}

func InitBundle() (*i18n.Bundle, error) {
	bundle := i18n.NewBundle(language.Chinese)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.MustLoadMessageFile("config/zh-CN.json")
	return bundle, nil
}

func RunIndex(ctx context.Context, opts ...Option) error {
	initLogger()
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

	if o.Index == "delete" {
		injector.IndexManager.DeleteAllIndex(ctx)
		return nil
	}
	injector.IndexManager.InitIndex(ctx)
	if o.Index == "load" {
		if len(config.C.Index.Languages) == 0 {
			return nil
		}
		indexName := injector.IndexManager.IndexName(config.C.Index.Languages[0]) // 只创建一个索引
		injector.IndexManager.LoadData(ctx, indexName, "")
		return nil
	}

	if o.UpdateDb > 0 {
		injector.HandlerManager.UpdateGroupInfo(o.UpdateDb)
	} else {
		go injector.HandlerManager.CheckDeleteMessage(ctx)
		go injector.GrpcServer.Run()
		injector.HandlerManager.Start(ctx, true)
	}
	return nil
}

func RunManager(ctx context.Context, opts ...Option) error {
	initLogger()
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

	go injector.HandlerManager.CheckDeleteMessage(ctx)
	injector.HandlerManager.Start(ctx, false)
	return nil
}

func RunGrpc(ctx context.Context, opts ...Option) error {
	initLogger()
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
	injector.IndexManager.InitIndex(ctx)

	return injector.GrpcServer.Run()
}

func initLogger() {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceQuote:      true,                  //键值对加引号
		TimestampFormat: "2006-01-02 15:04:05", //时间格式
		FullTimestamp:   true,
	})
}
