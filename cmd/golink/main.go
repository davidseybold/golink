package main

import (
	"log/slog"
	"os"

	"github.com/davidseybold/golink"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := golink.Run(logger); err != nil {
		logger.Error("failed to run", "error", err)
		os.Exit(1)
	}
}
