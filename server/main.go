package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/berquerant/pneutrinoutil/pkg/logx"
	"github.com/berquerant/pneutrinoutil/pkg/version"
	"github.com/berquerant/pneutrinoutil/server/alog"
	"github.com/berquerant/pneutrinoutil/server/config"
	_ "github.com/berquerant/pneutrinoutil/server/docs"
	"github.com/berquerant/pneutrinoutil/server/server"
	"github.com/spf13/pflag"
)

// @title pneutrinoutil API
// @version 1.0
// @description pneutrinoutil http server
// @host localhost:9101
// @basePath /v1
func main() {
	fs := pflag.NewFlagSet("main", pflag.ContinueOnError)
	fs.Usage = func() {
		fmt.Println(usage)
		fs.PrintDefaults()
	}
	fs.Bool("version", false, "print pneutrinoutil-server version")
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
		c.Close()
		os.Exit(1)
	}
	defer c.Close()

	alog.Setup(os.Stdout, c.SLogLevel())
	run(context.Background(), c)
}

func run(ctx context.Context, c *config.Config) {
	iCtx, stop := signal.NotifyContext(
		ctx,
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer stop()

	alog.L().Info("start server", c.LogAttr()...)
	srv := server.New(c)
	srv.Start(iCtx)
	srv.Wait()
	alog.L().Info("shut down")
}

const usage = `pneutrinoutil-server -- pneutrinoutil http server

e.g.
pneutrinoutil-server --neutrinoDir /path/to/NEUTRINO --workDir /path/to/workingDirectory --pneutrinoutil /path/to/pneutrinoutil

Flags:`
