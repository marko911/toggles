package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bugsnag/bugsnag-go"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	bugsnag.Configure(bugsnag.Configuration{
		// Your Bugsnag project API key
		APIKey: os.Getenv("BUGSNAG_API_KEY"),
		// The import paths for the Go packages containing your source files
		ProjectPackages: []string{"main", "toggles/server"},
	})

	app := &cli.App{
		Name:  "toggles server",
		Usage: "Feature flags api",
		Action: func(c *cli.Context) error {
			logrus.Println("Creating new server..")

			server := NewServer(c)
			logrus.Println("Server created, starting...")
			server.Start(c)
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "server-address",
				Usage:       "server address",
				EnvVars:     []string{"PORT"},
				DefaultText: "8080",
			},
			&cli.StringFlag{
				Name:        "database-address",
				Usage:       "Mongodb address",
				EnvVars:     []string{"DB_ADDRESS"},
				DefaultText: "localhost",
			},
			&cli.StringFlag{
				Name:        "database-name",
				Usage:       "Mongodb name",
				EnvVars:     []string{"DB_NAME"},
				DefaultText: "toggles",
			},
			&cli.StringFlag{
				Name:        "mongo-username",
				Usage:       "Mongodb login",
				EnvVars:     []string{"MONGO_USER"},
				DefaultText: "root",
			},
			&cli.StringFlag{
				Name:        "mongo-password",
				Usage:       "Mongodb pass",
				EnvVars:     []string{"MONGO_PASS"},
				DefaultText: "pass",
			},
			&cli.StringSliceFlag{
				EnvVars: []string{"SERVER_ALLOWED_HOSTS"},
				Name:    "server-allowed-hosts",
				Usage:   "server allowed hosts (CORS)",
			},
			&cli.StringSliceFlag{
				EnvVars: []string{"NATS_CLIENT_PASS"},
				Name:    "nats-client-pass",
				Usage:   "nats client password token",
			},
			&cli.StringSliceFlag{
				EnvVars: []string{"NATS_SERVER_URL"},
				Name:    "nats-server-url",
				Usage:   "nats url",
			},
			&cli.StringFlag{
				EnvVars: []string{"BUGSNAG_API_KEY"},
				Name:    "busnag-api-key",
				Usage:   "bugsnag",
			},
		},
	}
	fmt.Println("os args", os.Args)
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
