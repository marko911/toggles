package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "sc-server",
		Usage: "Serve graphql and auth requests",
		Action: func(c *cli.Context) error {
			fmt.Println("about to start server")
			server := NewServer(c)
			fmt.Println("new server made")
			server.Start(c)
			fmt.Println("server started")
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
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
