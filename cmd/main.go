package main

import (
	"context"
	"os"

	app "github.com/flytrap/telegram-bot/internal"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var VERSION = "0.2.9"

func main() {
	ctx := context.Background()
	app := cli.NewApp()
	app.Name = "telegram-bot"
	app.Version = VERSION
	app.Usage = "telegram-bot based on GRPC + WIRE."
	app.Commands = []*cli.Command{
		newIndexCmd(ctx), newGrpcCmd(ctx),
	}
	err := app.Run(os.Args)
	if err != nil {
		logrus.Error(err.Error())
	}
}

func newIndexCmd(ctx context.Context) *cli.Command {
	return &cli.Command{
		Name:  "index",
		Usage: "Run tg-bot index server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "conf",
				Aliases:  []string{"c"},
				Required: true,
				Usage:    "指定启动配置文件(.json,.yaml,.toml)",
			},
			&cli.IntFlag{
				Name:        "update",
				Aliases:     []string{"u"},
				Usage:       "更新索引原始数据",
				Required:    false,
				DefaultText: "0",
			},
			&cli.StringFlag{
				Name:        "index",
				Aliases:     []string{"i"},
				Usage:       "索引操作(load|delete)",
				Required:    false,
				DefaultText: "",
			},
		},
		Action: func(c *cli.Context) error {
			return app.RunIndex(ctx,
				app.SetConfigFile(c.String("conf")),
				app.SetVersion(VERSION),
				app.SetUpdateDb(c.Int64("update")),
				app.SetIndex(c.String("index")),
			)
		},
	}
}

func newGrpcCmd(ctx context.Context) *cli.Command {
	return &cli.Command{
		Name:  "grpc",
		Usage: "Run tg-bot grpc server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "conf",
				Aliases:  []string{"c"},
				Required: true,
				Usage:    "指定启动配置文件(.json,.yaml,.toml)",
			},
		},
		Action: func(c *cli.Context) error {
			return app.RunGrpc(ctx,
				app.SetConfigFile(c.String("conf")),
				app.SetVersion(VERSION),
			)
		},
	}
}
