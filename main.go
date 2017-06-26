package main

import (
	"os"
	"sort"

	"fmt"

	log "github.com/Sirupsen/logrus"
	cli "gopkg.in/urfave/cli.v1"
)

// NomadState ...
type NomadState struct {
	Info   map[string]string
	Groups []NomadTaskGroupState
}

// NomadTaskGroupState ...
type NomadTaskGroupState struct {
	Job   string
	Group string
	Count int
}

func main() {
	app := cli.NewApp()
	app.Name = "nomad-scale-helper"
	app.Usage = "easily restore / snapshot your nomad job scale config"
	app.Version = "0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "log-level",
			Value:  "info",
			Usage:  "Debug level (debug, info, warn/warning, error, fatal, panic)",
			EnvVar: "LOG_LEVEL",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "scale-export",
			Usage: "Export nomad job scale config to a local file",
			Action: func(c *cli.Context) error {
				configFile := c.Args().Get(0)
				if configFile == "" {
					return fmt.Errorf("Missing file name")
				}

				return exportCommand(configFile)
			},
		},
		{
			Name:  "scale-import",
			Usage: "Import nomad job scale config from a local file to Nomad cluster",
			Action: func(c *cli.Context) error {
				configFile := c.Args().Get(0)
				if configFile == "" {
					return fmt.Errorf("Missing file name")
				}

				return importCommand(configFile)
			},
		},
		{
			Name:  "drain",
			Usage: "Drain node and wait for all allocations to stop",
			Action: func(c *cli.Context) error {
				return drainCommand()
			},
		},
	}
	app.Before = func(c *cli.Context) error {
		// convert the human passed log level into logrus levels
		level, err := log.ParseLevel(c.String("log-level"))
		if err != nil {
			log.Fatal(err)
		}
		log.SetLevel(level)

		return nil
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	app.Run(os.Args)
}
