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
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	logger      *slog.Logger
	logLevel    *slog.LevelVar
	target      string
	matchLogger *slog.Logger

	filename  string        = "../lyrics.dat"
	delayWord time.Duration = 200 * time.Millisecond
	delayLine time.Duration = 1000 * time.Millisecond
)

func LyricsFromFile(filename string) (r *bufio.Reader, err error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return r, err
	}

	r = bufio.NewReader(bytes.NewBuffer(b))
	return r, nil
}

func processDelay(w string) bool {
	// Set up our command matching regexps
	//
	match := regexp.MustCompile(`{(?P<msec>\d+)}`)

	if !match.MatchString(w) {
		// This word is not a delay command, carry on
		//
		return false
	}

	// Pull the number from the matched string
	//
	msecString := match.FindStringSubmatch(w)[1]
	msec, err := strconv.ParseInt(msecString, 10, 64)
	if err != nil {
		logger.Info("unable to parse delay command, treating like a regular word", "word", w)
		return false
	}

	logger.Debug("delay command detected", "msec", msec)
	time.Sleep(time.Duration(msec) * time.Millisecond)

	return true
}

func newLine(o io.Writer) {
	fmt.Fprintf(o, "\n")
}

func singLine(o io.Writer, line string) error {
	if len(line) == 0 {
		newLine(o)
		return nil
	}

	// Sing a line of the lyrics with delays for embedded timing commands
	words := strings.Split(line, " ")

	for i, w := range words {
		addPadding, err := singWord(o, w)
		if err != nil {
			return err
		}
		if addPadding && i < (len(words)-1) {
			// That word calls for padding
			//
			fmt.Fprintf(o, " ")
		}
	}

	time.Sleep(delayLine)
	newLine(o)

	return nil
}

func singWord(o io.Writer, w string) (bool, error) {
	if processDelay(w) {
		// This word was a delay command, no further output
		return false, nil
	}

	_, err := fmt.Fprintf(o, "%s", w)
	if err != nil {
		return false, err
	}

	time.Sleep(delayWord)

	return true, nil
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

		err = singLine(o, string(line))
		if err != nil {
			return err
		}

		time.Sleep(delayLine)

	}
	return nil
}

func run(ctx context.Context, stdout io.Writer, stderr io.Writer, getenv func(string) string, args []string) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	setupLogger(ctx, stdout)

	logger.Info("Starting Rickroll",
		"filename", filename,
		"delayWord", delayWord,
		"delayLine", delayLine)

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

func cleanup() {
	logger.Info("Interrupt detected, exiting")
}

// main does as little as we can get away with.
func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()

	ctx := context.Background()

	if err := run(ctx, os.Stdout, os.Stderr, os.Getenv, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
