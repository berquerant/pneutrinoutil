package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"

	"github.com/berquerant/pneutrinoutil/pkg/logx"
	"github.com/berquerant/pneutrinoutil/pkg/wait"
	"github.com/berquerant/pneutrinoutil/server/alog"
	"github.com/berquerant/pneutrinoutil/server/config"
	"github.com/berquerant/pneutrinoutil/server/handler"
	"github.com/berquerant/pneutrinoutil/server/pworker"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	echoSwagger "github.com/swaggo/echo-swagger"
)

type Server struct {
	e      *echo.Echo
	c      *config.Config
	list   *pworker.List
	worker *wait.Worker
	wg     *sync.WaitGroup
}

func (s *Server) Echo() *echo.Echo { return s.e }
func (s *Server) Wait()            { s.wg.Wait() }

func New(cfg *config.Config) *Server {
	var (
		e      = echo.New()
		list   = pworker.NewList(cfg.PneutrinoutilWorkDir())
		worker = wait.New(cfg.Concurrency)
	)

	//
	// logger
	//
	l, ok := e.Logger.(*log.Logger)
	if !ok {
		panic(fmt.Sprintf("echo.Logger must be log.Logger"))
	}
	l.SetLevel(cfg.EchoLogLevel())
	//
	// middlewares
	//
	const healthPath = "/health"
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: cfg.AccessLogWriter,
		Skipper: func(c echo.Context) bool {
			// skip /health access log
			return strings.Contains(c.Request().RequestURI, healthPath)
		},
		Format: func() string {
			entries := []string{
				"time=${time_rfc3339_nano}",
				"id=${id}",
				"remote_ip=${remote_ip}",
				"host=${host}",
				"method=${method}",
				"uri=${uri}",
				"user_agent=${user_agent}",
				"status=${status}",
				"error=${error}",
				"latency=${latency}",
				"latency_human=${latency_human}",
				"bytes_in=${bytes_in}",
				"bytes_out=${bytes_out}",
			}
			return strings.Join(entries, "\t")
		}(),
		CustomTimeFormat: "2006-01-02 15:04:05.00000",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	//
	// handlers
	//
	v1 := e.Group("/v1")
	v1.GET(healthPath, handler.Health).Name = "health"
	v1.GET("/version", handler.Version).Name = "version"
	v1.GET("/debug", handler.Debug).Name = "debug"
	v1.GET("/swagger/*", echoSwagger.WrapHandler)
	v1.POST("/proc", handler.NewStart(cfg, worker).Handler).Name = "createProcess"
	v1.GET("/proc", handler.NewList(list).Handler).Name = "listProcess"
	getGroup := v1.Group("/proc/:id")
	getHandler := handler.NewGet(list, cfg.LogDir())
	getGroup.GET("/detail", getHandler.Detail).Name = "getDetail"
	getGroup.GET("/config", getHandler.Config).Name = "getConfig"
	getGroup.GET("/musicxml", getHandler.MusicXML).Name = "getMusicXML"
	getGroup.GET("/wav", getHandler.Wav).Name = "getWav"
	getGroup.GET("/world_wav", getHandler.WorldWav).Name = "getWorldWav"
	getGroup.GET("/log", getHandler.Log).Name = "getLog"
	//
	// echo
	//
	e.HideBanner = true
	e.HidePort = true

	var wg sync.WaitGroup
	return &Server{
		e:      e,
		c:      cfg,
		list:   list,
		worker: worker,
		wg:     &wg,
	}
}

func (s *Server) StartLoop(ctx context.Context) {
	if err := s.list.Init(); err != nil {
		alog.L().Error("init server", logx.Err(err))
	}

	//
	// wait process and write result loop
	//
	s.wg.Add(1)
	go func() {
		alog.L().Info("start worker loop")
		defer func() {
			s.wg.Done()
			alog.L().Info("end worker loop")
		}()

		for r := range s.worker.DoneC() {
			p, ok := r.Waiter.(*pworker.Process)
			if !ok {
				alog.L().Error("unknown waiter", slog.String("instance", fmt.Sprintf("%#v", r.Waiter)))
				continue
			}

			s.wg.Add(1)
			go func() {
				defer s.wg.Done()
				s.notify(ctx, p.RequestID, r.Err == nil)
			}()

			if err := r.Err; err != nil {
				alog.L().Error("waiter done", slog.String("id", p.RequestID), logx.Err(err))
				s.list.Fail(p.RequestID, err)
				continue
			}
			alog.L().Info("waiter done", slog.String("id", p.RequestID))
			if err := s.list.Init(); err != nil {
				alog.L().Error("load waiter", slog.String("id", p.RequestID), logx.Err(err))
			}
		}
	}()

	<-ctx.Done()
	alog.L().Info("server canceled, worker canceling")

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.worker.WaitAndClose()
	}()

	alog.L().Info("cancel worker")
	s.worker.Cancel()
}

func (s *Server) Start(ctx context.Context) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.StartLoop(ctx)
	}()

	go func() {
		alog.L().Info("start server")
		if err := s.e.Start(s.c.Addr()); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic("shutting down the server")
		}
	}()

	<-ctx.Done()
	alog.L().Info("server canceled")

	iCtx, cancel := context.WithTimeout(context.Background(), s.c.ShutdownPeriod())
	defer cancel()

	alog.L().Info("start shutting down")
	if err := s.e.Shutdown(iCtx); err != nil {
		panic(err)
	}
}
