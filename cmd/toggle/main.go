package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Name:  "sc-server",
		Usage: "Serve graphql and auth requests",
		Action: func(c *cli.Context) error {

			server := NewServer(c)
			server.Start(c)

			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "server-address",
				Usage:       "server address",
				EnvVars:     []string{"SERVER_ADDR"},
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
				EnvVars:     []string{"MONGO_INITDB_ROOT_USERNAME"},
				DefaultText: "root",
			},
			&cli.StringFlag{
				Name:        "mongo-password",
				Usage:       "Mongodb pass",
				EnvVars:     []string{"MONGO_INITDB_ROOT_PASSWORD"},
				DefaultText: "pass",
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
