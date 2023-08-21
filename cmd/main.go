package main

import (
	"context"
	"os"

	app "github.com/flytrap/telegram-bot/internal"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var VERSION = "0.1.0"

func main() {
	ctx := context.Background()
	app := cli.NewApp()
	app.Name = "telegram-bot"
	app.Version = VERSION
	app.Usage = "telegram-bot based on GRPC + WIRE."
	app.Commands = []*cli.Command{
		newWebCmd(ctx),
	}
	err := app.Run(os.Args)
	if err != nil {
		logrus.Error(err.Error())
	}
}

func newWebCmd(ctx context.Context) *cli.Command {
	return &cli.Command{
		Name:  "index",
		Usage: "Run tg-bot index server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "conf",
				Aliases: []string{"c"},
				Usage:   "App configuration file(.json,.yaml,.toml)",
			},
			&cli.IntFlag{
				Name:        "update",
				Aliases:     []string{"u"},
				Usage:       "更新索引数据",
				Required:    false,
				DefaultText: "0",
			},
		},
		Action: func(c *cli.Context) error {
			return app.RunIndex(ctx,
				app.SetConfigFile(c.String("conf")),
				app.SetVersion(VERSION),
				app.SetUpdateDb(c.Int("update")),
			)
		},
	}
}
