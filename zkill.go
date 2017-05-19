package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "zKillboard CLI"
	app.Usage = "I have no idea but guinness is a great beer."

	app.Commands = []cli.Command{{
		Name:    "listen",
		Aliases: []string{},
		Usage:   "view new kills as they occur",
		Action:  listener,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "alliance, a",
				Value:  "",
				Usage:  "User corporation/alliance",
				EnvVar: "ALLIANCE",
			},
			cli.BoolFlag{
				Name:  "filter-alliance",
				Usage: "Only show kills involving a corp/alliance",
			},
			cli.Float64Flag{
				Name:  "filter-isk",
				Value: 0.0,
				Usage: "Minimum threshhold to display a kill",
			},
			cli.StringFlag{
				Name:  "filter-ship",
				Value: "",
				Usage: "Only show specific ships",
			},
			cli.Float64Flag{
				Name:  "isk-threshhold, i",
				Value: 1000000000.0,
				Usage: "Highlight threshhold for isk amount",
			},
			cli.BoolFlag{
				Name:  "log",
				Usage: "Turns on logging of all kills",
			},
			cli.BoolTFlag{
				Name:  "verbose",
				Usage: "Listen output includes zkill URL",
			},
		},
	},
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "insecure",
			Usage:  "Skip SSL validation for corporate networks",
			Hidden: true,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err.Error())
		os.Exit(1)
	}
}
