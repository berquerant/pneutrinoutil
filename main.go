package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/berquerant/pneutrinoutil/pkg/cmd"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	cmd.Main(ctx)
}
