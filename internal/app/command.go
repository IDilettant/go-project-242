package app

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"code/pkg/code"
)

// New
func New() *cli.Command {
	return &cli.Command{
		Name:      "hexlet-path-size",
		Usage:     "print size of a file or directory",
		ArgsUsage: "<path>",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "human",
				Aliases: []string{"H"},
				Usage:   "human-readable sizes (auto-select unit)",
			},
			&cli.BoolFlag{
				Name:    "all",
				Aliases: []string{"a"},
				Usage:   "include hidden files and directories",
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() != 1 {
				_ = cli.ShowRootCommandHelp(cmd)

				return cli.Exit("\npath is required", 2)
			}

			path := cmd.Args().First()

			opts := code.Options{
				All: cmd.Bool("all"),
			}

			size, err := code.GetSize(path, opts)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			sizeStr := code.FormatSize(size, cmd.Bool("human"))
			fmt.Println(code.FormatOutput(sizeStr, path))

			return nil
		},
	}
}
