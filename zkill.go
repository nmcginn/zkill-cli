package main

import (
	"github.com/urfave/cli"
	"fmt"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "zKillboard CLI"
	app.Usage = "I have no idea but guinness is a great beer."

	app.Commands = []cli.Command {{
			Name:		"listen",
			Aliases:	[]string{},
			Usage:		"view new kills as they occur",
			Action:		listener,
			Flags:		[]cli.Flag {
				cli.Float64Flag{
					Name: "isk-threshhold, i",
					Value: 1000000000.0,
					Usage: "Highlight threshhold for isk amount",
				},
				cli.StringFlag{
					Name: "alliance, a",
					Value: "",
					Usage: "User corporation/alliance",
					EnvVar: "ALLIANCE",
				},
			},
		},
	}

	app.Flags = []cli.Flag {
		cli.BoolFlag{
			Name: "insecure",
			Usage: "Skip SSL validation for corporate networks",
			Hidden: true,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err.Error())
		os.Exit(1)
	}
}



