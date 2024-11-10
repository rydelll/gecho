package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/rydelll/gecho/cmd"
	"golang.org/x/sys/unix"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), unix.SIGINT, unix.SIGTERM, unix.SIGQUIT)
	defer cancel()

	if err := cmd.Execute(ctx, os.Args, os.Getenv, os.Stderr); err != nil {
		panic(err)
	}
}
