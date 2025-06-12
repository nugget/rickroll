package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
)

var (
	filename  string = "../lyrics.dat"
	delayWord int    = 200
	delayLine int    = 1000

	logger      *slog.Logger
	logLevel    *slog.LevelVar
	target      string
	matchLogger *slog.Logger
)

func singLine(line string) error {
	fmt.Println("# LINE", line)

	fmt.Printf("\n")
	return nil
}

func setupLogger(ctx context.Context, stdout io.Writer) {
	logLevel = new(slog.LevelVar)

	handlerOptions := &slog.HandlerOptions{
		Level: logLevel,
	}
	handler := slog.NewTextHandler(stdout, handlerOptions)
	logger = slog.New(handler)
}

func run(ctx context.Context, stdout io.Writer, stderr io.Writer, getenv func(string) string, args []string) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	setupLogger(ctx, stdout)

	b, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	r := bufio.NewReader(bytes.NewBuffer(b))
	for {
		line, prefix, err := r.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			logger.Error("error reading line",
				"line", line,
				"prefix", prefix,
				"error", err)
			return err
		}

		if len(line) == 0 {
			fmt.Printf("\n")
		} else {
			err = singLine(string(line))
			if err != nil {
				return err
			}
		}

	}
	return nil
}

// main does as little as we can get away with.
func main() {
	ctx := context.Background()

	if err := run(ctx, os.Stdout, os.Stderr, os.Getenv, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
