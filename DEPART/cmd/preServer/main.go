package main

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"os"
	"path/filepath"
	"preServer/conf"
	"preServer/server"
	"preServer/utils"
	"runtime"
)

var (
	app *cli.App
)

func init() {
	app = NewApp()
	app.Flags = []cli.Flag{
		utils.HttpPortFlag,
	}
	app.Action = startServer
	app.Before = beforeStart
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = filepath.Base(os.Args[0])
	app.Author = "wxblockchain"
	app.Email = ""
	app.Version = "1.0.0"
	app.Usage = "Pre server"
	return app
}

func beforeStart(ctx *cli.Context) error {
	runtime.GOMAXPROCS(runtime.NumCPU())
	return nil
}

func startServer(ctx *cli.Context) error {
	cfg, err := conf.MakeConfig(ctx)
	if err != nil {
		return err
	}

	serv, err := server.NewServer(cfg)
	if err != nil {
		return err
	}
	err = serv.Start()
	if err != nil {
		return err
	}
	serv.Wait()
	return nil
}
