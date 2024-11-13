package cmd

import (
	"context"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/rydelll/gecho/internal/logging"
	"github.com/rydelll/gecho/internal/server"
)

// defaultPort sets the default port for a server to listen on.
const defaultPort = 7777

// Execute parses arguments and environment variables, initializes dependencies,
// and starts the application.
func Execute(ctx context.Context, args []string, env func(string) string, stderr io.Writer) error {
	// Arguments
	fs := flag.NewFlagSet("", flag.ExitOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprintf(stderr, "Gecho Server\n\n")
		fmt.Fprintf(stderr, "Usage:\n\n")
		fmt.Fprintf(stderr, "\t%s [options]\n\n", args[0])
		fmt.Fprintf(stderr, "Options:\n")
		fs.PrintDefaults()
		fmt.Fprintln(stderr)
	}

	var port uint
	fs.UintVar(&port, "port", defaultPort, "port for the server to listen on")
	fs.Parse(args[1:])

	// Logger
	logLevel := logging.SlogLevel(env("LOG_LEVEL"))
	logJSON := strings.ToLower(env("LOG_MODE")) == "json"
	logger := logging.NewLogger(stderr, logLevel, logJSON)

	// Server
	srv, err := server.New(logger, port)
	if err != nil {
		return fmt.Errorf("init server: %w", err)
	}

	srv.ListenAndServe(ctx)

	return nil
}
