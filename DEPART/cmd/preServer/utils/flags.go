package utils

import "gopkg.in/urfave/cli.v1"

var HttpPortFlag = cli.StringFlag{
	Name:  "port",
	Usage: "http port",
}

