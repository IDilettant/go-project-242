package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:  "hexlet-path-size",
		Usage: "print size of a file or directory",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			fmt.Println("Hello from Hexlet!")

			return nil
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		os.Exit(1)
	}
}
