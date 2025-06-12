package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/nugget/go-telnet"
)

var (
	logger      *slog.Logger
	logLevel    *slog.LevelVar
	target      string
	matchLogger *slog.Logger

	filename  string        = "lyrics.dat"
	delayWord time.Duration = 200 * time.Millisecond
	delayLine time.Duration = 1000 * time.Millisecond

	lyricsBytes []byte
)

func LoadLyricsFromFile(filename string) (err error) {
	lyricsBytes, err = os.ReadFile(filename)
	if err != nil {
		return err
	}

	logger.Info("loaded lyrics from file", "file", filename, "bytes", len(lyricsBytes))

	return nil
}

func LyricsReader() (r *bufio.Reader, err error) {
	r = bufio.NewReader(bytes.NewBuffer(lyricsBytes))

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
	fmt.Fprintf(o, "\r\n")
}

func singLine(o io.Writer, line string) error {
	if len(line) == 0 {
		newLine(o)
		return nil
	}

	logger.Debug("singing line", "line", line)

	// Sing a line of the lyrics with delays for embedded timing commands
	words := strings.Split(line, " ")

	wordCount := len(words) - 1

	c := words[wordCount][0]
	if c == 123 {
		// The last "word" of the line is a delay command, so we
		// don't want to pad with a space after the second-to-last
		// word (which is really the last word)
		wordCount = wordCount - 1
	}

	for i, w := range words {
		syllables := strings.Split(w, "-")
		if len(syllables) > 1 {
			logger.Debug("Syllables detected!", "syllables", syllables)

			for _, s := range syllables {
				singWord(o, s)
			}

			if i < (len(words) - 1) {
				fmt.Fprintf(o, " ")
			}
		} else {
			addPadding, err := singWord(o, w)
			if err != nil {
				return err
			}
			logger.Debug("end of word", "w", w, "addPadding", addPadding, "i", i, "wordCount", wordCount)
			if addPadding && i < wordCount {
				// That word calls for padding
				//
				fmt.Fprintf(o, " ")
			}
		}
	}

	time.Sleep(delayLine)
	newLine(o)

	return nil
}

func singWord(o io.Writer, w string) (bool, error) {
	logger.Debug("singing word", "w", w)
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

func SessionHandler(ctx context.Context, stdout io.Writer) error {
	r, err := LyricsReader()
	if err != nil {
		return err
	}

	err = SingSong(r, stdout)
	if err != nil {
		return err
	}

	return nil
}

type internalTelnetHandler struct{}

var TelnetHandler telnet.Handler = internalTelnetHandler{}

func dnsLookup(remoteString string) (remoteAddr, remotePort, remoteHost string) {
	var err error

	remoteAddr, remotePort, err = net.SplitHostPort(remoteString)
	if err != nil {
		logger.Info("unable to split host/port from remoteAddr", "error", err, "remoteString", remoteString)
	} else {
		var lookupSlice []string

		lookupSlice, err := net.LookupAddr(remoteAddr)
		if err != nil {
			logger.Debug("reverse dns for client not found", "remoteAddr", remoteAddr, "error", err)
		}
		if len(lookupSlice) > 0 {
			remoteHost = lookupSlice[0]
		} else {
			logger.Debug("no reverse dns found", "remoteAddr", remoteAddr, "lookupSlice", lookupSlice, "length", len(lookupSlice))
		}
	}

	return remoteAddr, remotePort, remoteHost
}

func (handler internalTelnetHandler) ServeTELNET(ctx telnet.Context, w telnet.Writer, r telnet.Reader) {
	var err error

	conn := ctx.Conn()
	defer conn.Close()

	remoteString := conn.RemoteAddr().String()

	remoteAddr, remotePort, remoteHost := dnsLookup(remoteString)

	startTime := time.Now()

	logger.Info("connection opened", "remoteHost", remoteHost, "remoteAddr", remoteAddr, "remotePort", remotePort)

	c := context.Background()
	err = SessionHandler(c, w)
	duration := time.Now().Sub(startTime)
	if err != nil {
		logger.Info("connection error", "error", err, "remoteHost", remoteHost, "remoteAddr", remoteAddr, "remotePort", remotePort, "duration", duration)
		return
	}

	err = conn.Close()
	if err != nil {
		logger.Info("connection lost", "error", err, "remoteHost", remoteHost, "remoteAddr", remoteAddr, "remotePort", remotePort, "duration", duration)
	} else {
		logger.Info("connection closed", "remoteHost", remoteHost, "remoteAddr", remoteAddr, "remotePort", remotePort, "duration", duration)
	}
}

func LaunchTelnetServer(ctx context.Context) error {
	logger.Info("starting telnet server")

	err := telnet.ListenAndServe(":23", TelnetHandler)
	if err != nil {
		return err
	}

	return nil
}

func run(ctx context.Context, stdout io.Writer, stderr io.Writer, getenv func(string) string, args []string) (err error) {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	setupLogger(ctx, stdout)

	commit, buildDate, _ := getBuildInfo()

	logger.Info("starting rickrolld",
		"filename", filename,
		"delayWord", delayWord,
		"delayLine", delayLine,
		"commit", commit,
		"buildDate", buildDate)

	err = LoadLyricsFromFile(filename)
	if err != nil {
		return err
	}

	err = LaunchTelnetServer(ctx)
	if err != nil {
		return err
	}

	logger.Info("Done")

	return nil
}

func getBuildInfo() (commit, buildDate string, dirty bool) {
	buildInfo, ok := debug.ReadBuildInfo()

	if !ok {
		return
	}

	dirty = false

	for _, setting := range buildInfo.Settings {
		switch setting.Key {
		case "vcs.revision":
			commit = setting.Value
		case "vcs.time":
			buildDate = setting.Value
		case "vcs.modified":
			dirty = true
		}
	}

	return
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
