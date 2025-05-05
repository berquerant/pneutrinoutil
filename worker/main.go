package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/berquerant/pneutrinoutil/pkg/alog"
	"github.com/berquerant/pneutrinoutil/pkg/logx"
	"github.com/berquerant/pneutrinoutil/pkg/version"
	"github.com/berquerant/pneutrinoutil/worker/config"
	"github.com/berquerant/pneutrinoutil/worker/worker"
	"github.com/spf13/pflag"
)

func main() {
	fs := pflag.NewFlagSet("main", pflag.ContinueOnError)
	fs.Usage = func() {
		fmt.Println(usage)
		fs.PrintDefaults()
	}

	fs.Bool("version", false, "print pneutrinoutil-worker version")
	c, err := config.New(fs)
	if errors.Is(err, pflag.ErrHelp) {
		return
	}
	if err != nil {
		panic(err)
	}
	if v, _ := fs.GetBool("version"); v {
		version.Write(os.Stdout)
		return
	}

	err = c.Init()
	if err != nil {
		alog.L().Error("invalid config", logx.Err(err))
		_ = c.Close()
		os.Exit(1)
	}
	defer func() { _ = c.Close() }()

	alog.Setup(os.Stdout, c.SLogLevel())
	run(context.Background(), c)
}

func run(ctx context.Context, c *config.Config) {
	iCtx, stop := signal.NotifyContext(
		ctx,
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer stop()

	srv := worker.New(c)
	alog.L().Info("start worker")
	if err := srv.Run(iCtx); err != nil {
		alog.L().Error("worker got an error", logx.Err(err))
	}
	alog.L().Info("shut down")
}

const usage = `pneutrinoutil-worker -- pneutrinoutil process worker

e.g.
pneutrinoutil-worker --neutrinoDir /path/to/NEUTRINO --workDir /path/to/workingDirectory --pneutrinoutil /path/to/pneutrinoutil --mysqlDSN DSN --redisDSN DSN

Flags:`
