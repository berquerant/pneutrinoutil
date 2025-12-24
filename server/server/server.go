package server

import (
	"context"
	"database/sql"
	"errors"
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
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	echoSwagger "github.com/swaggo/echo-swagger"
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
	// logger
	//
	l, ok := e.Logger.(*log.Logger)
	if !ok {
		panic("echo.Logger must be log.Logger")
	}
	l.SetLevel(cfg.EchoLogLevel())
	//
	// middlewares
	//
	const healthPath = "/health"
	e.Use(middleware.RequestLoggerWithConfig(middleware.LoggerConfig{
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
			}
			if cfg.Debug {
				entries = append(entries,
					"bytes_in=${bytes_in}",
					"bytes_out=${bytes_out}",
				)
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
	v1.POST("/proc", handler.NewStart(client, cfg.ProcessTimeout(), cfg.StorageBucket, cfg.StoragePath, objectAdmin, details, processes).Handler).Name = "createProcess"
	v1.GET("/proc/search", handler.NewSearch(searcher).SearchProcess).Name = "searchProcess"
	getGroup := v1.Group("/proc/:id")
	getHandler := handler.NewGet(processes, details, objectAdmin, objects)
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
