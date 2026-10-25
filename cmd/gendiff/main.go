package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"code"

	"github.com/urfave/cli/v3"
)

func main() {
	if err := newCommand(os.Stdout, os.Stderr).Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newCommand(output, errorOutput io.Writer) *cli.Command {
	return &cli.Command{
		Name:      "gendiff",
		Usage:     "Compares two configuration files and shows a difference.",
		ArgsUsage: "<filepath1> <filepath2>",
		Writer:    output,
		ErrWriter: errorOutput,
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "format", Aliases: []string{"f"}, Value: "stylish", Usage: "output format"},
		},
		Action: func(_ context.Context, command *cli.Command) error {
			if command.NArg() != 2 {
				return errors.New("expected two file paths; use gendiff --help for usage")
			}
			result, err := code.GenDiff(command.Args().Get(0), command.Args().Get(1), command.String("format"))
			if err != nil {
				return err
			}
			if result == "" {
				return nil
			}
			_, err = fmt.Fprintln(command.Writer, result)
			return err
		},
	}
}
