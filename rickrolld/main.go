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

func LyricsFromFile(filename string) (r *bufio.Reader, err error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return r, err
	}

	r = bufio.NewReader(bytes.NewBuffer(b))
	return r, nil
}

func singLine(o io.Writer, line string) error {
	fmt.Fprintln(o, "# LINE", line)
	fmt.Fprintf(o, "\n")
	return nil
}

func SingSong(r *bufio.Reader, o io.Writer) error {
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
			fmt.Fprintf(o, "\n")
		} else {
			err = singLine(o, string(line))
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func run(ctx context.Context, stdout io.Writer, stderr io.Writer, getenv func(string) string, args []string) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	setupLogger(ctx, stdout)
	logger.Info("Starting Rickroll")

	r, err := LyricsFromFile(filename)
	if err != nil {
		return err
	}

	err = SingSong(r, stdout)
	if err != nil {
		return err
	}

	logger.Info("Done")

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

// main does as little as we can get away with.
func main() {
	ctx := context.Background()

	if err := run(ctx, os.Stdout, os.Stderr, os.Getenv, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
