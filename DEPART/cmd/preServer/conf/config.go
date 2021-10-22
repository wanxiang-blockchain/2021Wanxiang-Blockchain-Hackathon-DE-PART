package conf

import (
	"gopkg.in/urfave/cli.v1"
	"preServer/utils"
)

type ServConfig struct {
	HttpCfg HttpConfig
}

type HttpConfig struct {
	Ip   string
	Port string
}

func defaultConfig() *ServConfig {
	return &ServConfig{
		HttpCfg: HttpConfig{
			Ip:   "127.0.0.1",
			Port: "8080",
		},
	}
}

func MakeConfig(ctx *cli.Context) (*ServConfig, error) {
	config := defaultConfig()
	port := ctx.GlobalString(utils.HttpPortFlag.Name)
	if port != "" {
		config.HttpCfg.Port = port
	}

	return config, nil
}
