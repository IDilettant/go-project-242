package app

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"code"
)

// New application constructor
func New() *cli.Command {
	return &cli.Command{
		Name:      "hexlet-path-size",
		Usage:     "print size of a file or directory; supports -r (recursive), -H (human-readable), -a (include hidden)",
		ArgsUsage: "<path>",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "human",
				Aliases:     []string{"H"},
				Usage:       "human-readable sizes (auto-select unit)",
				DefaultText: "false",
			},
			&cli.BoolFlag{
				Name:        "all",
				Aliases:     []string{"a"},
				Usage:       "include hidden files and directories",
				DefaultText: "false",
			},
			&cli.BoolFlag{
				Name:        "recursive",
				Aliases:     []string{"r"},
				Usage:       "recursive size of directories",
				DefaultText: "false",
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() != 1 {
				_ = cli.ShowRootCommandHelp(cmd)

				return cli.Exit("\npath is required", 2)
			}

			path := cmd.Args().First()

			output, err := code.GetPathSize(path, cmd.Bool("recursive"), cmd.Bool("human"), cmd.Bool("all"))
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			fmt.Println(output)

			return nil
		},
	}
}
