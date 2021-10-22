package main

import (
	"os"

	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "hospitalClient"
	app.Version = "1.0.0"

	app.Commands = []cli.Command{
		{
			Name:   "upload",
			Action: handleUpload,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "data",
					Value: "data",
					Usage: "patient's medical record",
				}, cli.StringFlag{
					Name:  "patientId",
					Value: "patientId",
					Usage: "patient's ID",
				}, cli.StringFlag{
					Name:  "patientPK",
					Usage: "patient public key",
					Value: "patientPK",
				},
			},
		},
		{
			Name:   "download",
			Usage:  "download patient record and check",
			Action: handleDownload,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "myPrivateKey",
					Value: "myPrivateKey",
					Usage: "my private key",
				}, cli.StringFlag{
					Name:  "rk",
					Value: "rk",
					Usage: "re encryption rk",
				}, cli.StringFlag{
					Name:  "pubX",
					Value: "pubX",
					Usage: "re encryption pubX",
				}, cli.StringFlag{
					Name:  "patientId",
					Value: "patientId",
					Usage: "patient's ID",
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
