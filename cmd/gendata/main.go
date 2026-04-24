package main

import (
	"log/slog"
	"os"

	"gendata/pkg/cmd"
)

func main() {
	cmd, err := cmd.NewRootCmd(os.Args[1:])
	if err != nil {
		slog.Error("New Command Failed", slog.Any("error", err))
		os.Exit(1)
	}
	if err := cmd.Execute(); err != nil {
		slog.Error("Execute Command Failed", slog.Any("error", err))
		os.Exit(1)
	}
}
