package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/berquerant/pneutrinoutil/pkg/alog"
	"github.com/berquerant/pneutrinoutil/pkg/domain"
	"github.com/berquerant/pneutrinoutil/pkg/infra"
	"github.com/berquerant/pneutrinoutil/pkg/logx"
	"github.com/berquerant/pneutrinoutil/pkg/repo"
	"github.com/berquerant/pneutrinoutil/server/config"
	"github.com/berquerant/pneutrinoutil/server/handler"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	echoSwagger "github.com/swaggo/echo-swagger/v2"
)

type Server struct {
	e      *echo.Echo
	c      *config.Config
	db     *sql.DB
	client *asynq.Client
}

func (s *Server) Echo() *echo.Echo { return s.e }

func (s *Server) Close() error {
	return errors.Join(s.db.Close(), s.client.Close())
}

func New(ctx context.Context, cfg *config.Config) (*Server, error) {
	e := echo.New()
	db, err := cfg.NewSQL(ctx)
	if err != nil {
		return nil, err
	}

	objectStorage, err := cfg.NewStorage(ctx)
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	client, err := cfg.NewAsynqClient()
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	var (
		objectConn   = infra.NewConn[domain.Object](db)
		objects      = repo.NewObject(objectConn, objectConn)
		objectAdmin  = repo.NewObjectAdmin(objects, objects, objectStorage, objectStorage)
		detailConn   = infra.NewConn[domain.ProcessDetails](db)
		details      = repo.NewProcessDetails(detailConn, detailConn)
		processConn  = infra.NewConn[domain.Process](db)
		processes    = repo.NewProcess(processConn, processConn)
		searcherConn = infra.NewConn[repo.SearchProcessResultElement](db)
		searcher     = repo.NewSearcher(searcherConn)
	)

	//
	// middlewares
	//
	const healthPath = "/health"
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		Skipper: func(c *echo.Context) bool {
			// skip /health access log
			return strings.Contains(c.Request().RequestURI, healthPath)
		},
		LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
			type pair struct {
				k string
				v any
			}
			xs := []pair{
				{k: "time", v: v.StartTime.Format("2006-01-02T15:04:05.00000")},
				{k: "id", v: v.RequestID},
				{k: "remote_ip", v: v.RemoteIP},
				{k: "host", v: v.Host},
				{k: "method", v: v.Method},
				{k: "uri", v: v.URI},
				{k: "referer", v: v.Referer},
				{k: "user_agent", v: v.UserAgent},
				{k: "status", v: v.Status},
				{k: "error", v: v.Error},
				{k: "latency", v: v.Latency},
				{k: "bytes_in", v: v.ContentLength},
				{k: "bytes_out", v: v.ResponseSize},
			}
			ys := make([]string, len(xs))
			for i, x := range xs {
				ys[i] = fmt.Sprintf("%s=%v", x.k, x.v)
			}
			_, err := fmt.Fprintln(cfg.AccessLogWriter, strings.Join(ys, "\t"))
			return err
		},
		LogLatency:       true,
		LogProtocol:      true,
		LogRemoteIP:      true,
		LogHost:          true,
		LogMethod:        true,
		LogURI:           true,
		LogRequestID:     true,
		LogReferer:       true,
		LogUserAgent:     true,
		LogStatus:        true,
		LogContentLength: true,
		LogResponseSize:  true,
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	//
	// handlers
	//
	v1 := e.Group("/v1")
	r1 := v1.GET(healthPath, handler.Health)
	r1.Name = "health"
	r2 := v1.GET("/version", handler.Version)
	r2.Name = "version"
	r3 := v1.GET("/debug", handler.Debug)
	r3.Name = "debug"
	v1.GET("/swagger/*", echoSwagger.WrapHandler)
	r4 := v1.POST("/proc", handler.NewStart(client, cfg.ProcessTimeout(), cfg.StorageBucket, cfg.StoragePath, objectAdmin, details, processes).Handler)
	r4.Name = "createProcess"
	r5 := v1.GET("/proc/search", handler.NewSearch(searcher).SearchProcess)
	r5.Name = "searchProcess"
	getGroup := v1.Group("/proc/:id")
	getHandler := handler.NewGet(processes, details, objectAdmin, objects)
	r6 := getGroup.GET("/detail", getHandler.Detail)
	r6.Name = "getDetail"
	r7 := getGroup.GET("/config", getHandler.Config)
	r7.Name = "getConfig"
	r8 := getGroup.GET("/musicxml", getHandler.MusicXML)
	r8.Name = "getMusicXML"
	r9 := getGroup.GET("/wav", getHandler.Wav)
	r9.Name = "getWav"
	r10 := getGroup.GET("/log", getHandler.Log)
	r10.Name = "getLog"

	return &Server{
		e:      e,
		c:      cfg,
		db:     db,
		client: client,
	}, nil
}

func (s *Server) Start(ctx context.Context) {
	defer func() {
		if err := s.Close(); err != nil {
			alog.L().Error("close server", logx.Err(err))
		}
	}()
	srv := &http.Server{
		Addr:    s.c.Addr(),
		Handler: s.e,
	}
	go func() {
		alog.L().Info("start server")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic("shutting down the server")
		}
	}()

	<-ctx.Done()
	alog.L().Info("server canceled")

	iCtx, cancel := context.WithTimeout(context.Background(), s.c.ShutdownPeriod())
	defer cancel()

	alog.L().Info("start shutting down")
	if err := srv.Shutdown(iCtx); err != nil {
		panic(err)
	}
}
