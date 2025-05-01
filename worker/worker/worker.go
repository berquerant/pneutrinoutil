package worker

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/berquerant/pneutrinoutil/pkg/alog"
	"github.com/berquerant/pneutrinoutil/pkg/domain"
	"github.com/berquerant/pneutrinoutil/pkg/infra"
	"github.com/berquerant/pneutrinoutil/pkg/logx"
	"github.com/berquerant/pneutrinoutil/pkg/repo"
	"github.com/berquerant/pneutrinoutil/pkg/task"
	"github.com/berquerant/pneutrinoutil/worker/config"
	"github.com/hibiken/asynq"
)

func New(c *config.Config) *Server {
	return &Server{
		c: c,
	}
}

type Server struct {
	c              *config.Config
	srv            *asynq.Server
	db             *sql.DB
	objects        *repo.ObjectAdmin
	processDetails *repo.ProcessDetails
	processes      *repo.Process
	webhook        *infra.Webhook
}

func (s *Server) Run(ctx context.Context) error {
	alog.L().Info("init worker")
	if err := s.init(ctx); err != nil {
		return err
	}
	defer func() {
		if err := s.close(); err != nil {
			alog.L().Error("closing worker", logx.Err(err))
		}
		alog.L().Info("close worker")
	}()
	alog.L().Info("start worker")
	return s.run()
}

func (s *Server) init(ctx context.Context) error {
	db, err := s.c.NewSQL(ctx)
	if err != nil {
		return err
	}
	s.db = db

	// TODO: use object storage
	storageObjects := infra.NewFileSystem(s.c.StorageDir)
	objectConn := infra.NewConn[domain.Object](db)
	objectTable := repo.NewObject(objectConn, objectConn)
	s.objects = repo.NewObjectAdmin(objectTable, objectTable, storageObjects, storageObjects)

	processDetailsConn := infra.NewConn[domain.ProcessDetails](db)
	s.processDetails = repo.NewProcessDetails(processDetailsConn, processDetailsConn)

	processConn := infra.NewConn[domain.Process](db)
	s.processes = repo.NewProcess(processConn, processConn)

	if x := s.c.NewWebhook(); x != nil {
		s.webhook = x
	}

	redisOpt, err := asynq.ParseRedisURI(s.c.RedisDSN)
	if err != nil {
		return err
	}
	s.srv = asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: s.c.Concurrency,
			Queues: map[string]int{
				"default": 10,
			},
			ShutdownTimeout: s.c.ShutdownPeriod(),
			Logger:          NewAsynqLogger(alog.L()),
			LogLevel:        s.c.AsynqLogLevel(),
		},
	)

	return nil
}

func (s *Server) newServeMux() *asynq.ServeMux {
	mux := asynq.NewServeMux()
	pneutrinoutilProcessor := task.NewPneutrinoutilProcessor(&task.PneutrinoutilProcessorParams{
		Pneutrinoutil: s.c.Pneutrinoutil,
		NeutrinoDir:   s.c.NeutrinoDir,
		WorkDir:       filepath.Join(s.c.WorkDir, "pneutrinoutil-processor"),
		Shell:         s.c.Shell,
		Bucket:        s.c.StorageBucket,
		BasePath:      s.c.StoragePath,
		Env: []string{
			fmt.Sprintf("HOME=%s", os.Getenv("HOME")),
			fmt.Sprintf("PWD=%s", os.Getenv("PWD")),
		},
		Webhooker:             s.webhook,
		ObjectReader:          s.objects,
		ObjectWriter:          s.objects,
		ProcessDetailsGetter:  s.processDetails,
		ProcessDetailsUpdater: s.processDetails,
		ProcessGetter:         s.processes,
		ProcessUpdater:        s.processes,
	})
	mux.HandleFunc(task.TypePneutrinoutilStart, pneutrinoutilProcessor.ProcessStart)
	return mux
}

func (s *Server) run() error {
	mux := s.newServeMux()
	return s.srv.Run(mux)
}

func (s *Server) close() error {
	return s.db.Close()
}
