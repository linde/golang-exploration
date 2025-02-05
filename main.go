package main

import (
	"log/slog"
	"myapp/cmd"
	"os"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cmd.Execute()
}
