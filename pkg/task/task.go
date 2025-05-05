package task

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"al.essio.dev/pkg/shellescape"
	"github.com/berquerant/pneutrinoutil/pkg/alog"
	"github.com/berquerant/pneutrinoutil/pkg/domain"
	"github.com/berquerant/pneutrinoutil/pkg/infra"
	"github.com/berquerant/pneutrinoutil/pkg/logx"
	"github.com/berquerant/pneutrinoutil/pkg/pathx"
	"github.com/berquerant/pneutrinoutil/pkg/ptr"
	"github.com/berquerant/pneutrinoutil/pkg/repo"
	"github.com/hibiken/asynq"
)

const (
	TypePneutrinoutilStart = "pneutrinoutil:start"
)

type PneutrinoutilStartPayload struct {
	RequestID string   `json:"rid"`
	Args      []string `json:"args"`
}

func NewPneutrinoutilStart(p PneutrinoutilStartPayload) (*asynq.Task, error) {
	payload, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypePneutrinoutilStart, payload), nil
}

type PneutrinoutilStartWebhookParams struct {
	Type      string `json:"type"`
	RequestID string `json:"rid"`
	OK        bool   `json:"ok"`
}

type PneutrinoutilProcessorParams struct {
	Pneutrinoutil string
	NeutrinoDir   string
	WorkDir       string
	Shell         string
	Bucket        string
	BasePath      string
	Env           []string

	Webhooker             infra.Webhooker // optional
	ObjectReader          repo.ObjectReader
	ObjectWriter          repo.ObjectWriter
	ProcessDetailsGetter  repo.ProcessDetailsGetter
	ProcessDetailsUpdater repo.ProcessDetailsUpdater
	ProcessGetter         repo.ProcessGetter
	ProcessUpdater        repo.ProcessUpdater
}

func NewPneutrinoutilProcessor(params *PneutrinoutilProcessorParams) *PneutrinoutilProcessor {
	return &PneutrinoutilProcessor{
		params,
	}
}

type PneutrinoutilProcessor struct {
	*PneutrinoutilProcessorParams
}

var (
	ErrTask = errors.New("Task")
)

func (p *PneutrinoutilProcessor) ProcessStart(ctx context.Context, t *asynq.Task) error {
	var (
		logAttrs = []any{"type", t.Type()}
		attrs    = func(v ...any) []any { return append(logAttrs, v...) }
		baseErr  = fmt.Errorf("%w: %s", ErrTask, t.Type())
	)

	alog.L().Info("got task", attrs()...)
	alog.L().Debug("got task", attrs("payload", string(t.Payload()))...)
	var payload PneutrinoutilStartPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return errors.Join(baseErr, err)
	}

	logAttrs = append(logAttrs, "rid", payload.RequestID)
	baseErr = fmt.Errorf("%w: rid=%s", baseErr, payload.RequestID)
	withBaseErr := func(err error, format string, v ...any) error {
		verr := fmt.Errorf("%w: %s", errors.Join(baseErr, err), fmt.Sprintf(format, v...))
		alog.L().Error("got error", attrs(logx.Err(verr))...)
		return verr
	}

	alog.L().Info("create work dir", attrs("dir", p.WorkDir)...)
	workDir := filepath.Join(p.WorkDir, payload.RequestID)
	if err := pathx.EnsureDir(workDir); err != nil {
		return withBaseErr(err, "failed to create work dir")
	}

	alog.L().Info("get process", attrs()...)
	proc, err := p.ProcessGetter.GetProcessByRequestId(ctx, payload.RequestID)
	if err != nil {
		return withBaseErr(err, "failed to get process")
	}
	if proc.Status != domain.ProcessStatusPending {
		return fmt.Errorf("%w: status is not pending", baseErr)
	}
	if _, err := p.ProcessUpdater.UpdateProcess(ctx, &repo.UpdateProcessRequest{
		ID:        proc.ID,
		Status:    ptr.To(domain.ProcessStatusRunning),
		StartedAt: ptr.To(time.Now()),
	}); err != nil {
		return withBaseErr(err, "failed to update process(%d)", proc.ID)
	}

	var processSucceed bool
	defer func() {
		if err := p.updateProcessStatus(ctx, proc.ID, time.Now(), processSucceed); err != nil {
			alog.L().Error("update process status", attrs("succeed", processSucceed, logx.Err(err))...)
		} else {
			alog.L().Info("update process status", attrs("succeed", processSucceed)...)
		}
		if err := p.webhook(ctx, payload.RequestID, processSucceed); err != nil {
			alog.L().Error("webhook failed", attrs(logx.Err(err))...)
		}
	}()

	alog.L().Info("get process details", attrs("id", proc.DetailsID)...)
	details, err := p.ProcessDetailsGetter.GetProcessDetails(ctx, proc.DetailsID)
	if err != nil {
		return withBaseErr(err, "process_details(%d) not found", proc.DetailsID)
	}
	alog.L().Info("get score object", attrs("id", details.ScoreObjectID)...)
	scoreObject, err := p.ObjectReader.ReadObject(ctx, details.ScoreObjectID)
	if err != nil {
		return withBaseErr(err, "score object(%d) not found", details.ScoreObjectID)
	}
	if scoreObject.Object().Type != domain.ObjectTypeFile {
		return fmt.Errorf("%w: score object(%d) is not a file", baseErr, details.ScoreObjectID)
	}

	score, _ := scoreObject.Storage()
	scorePath := filepath.Join(workDir, filepath.Base(score.Path))
	alog.L().Info("create local score file", attrs("path", scorePath)...)
	if err := func() error {
		f, err := os.Create(scorePath)
		if err != nil {
			return err
		}
		defer func() { _ = f.Close() }()
		_, err = io.Copy(f, score.Blob)
		return err
	}(); err != nil {
		return withBaseErr(err, "create local score file")
	}

	logPath := filepath.Join(workDir, "process.log")
	alog.L().Info("create local log file", attrs("path", logPath)...)
	logFile, err := os.Create(logPath)
	if err != nil {
		return withBaseErr(err, "failed to create local log file")
	}

	args := p.generateArgs(&payload, workDir, scorePath)
	alog.L().Info("start pneutrinoutil", attrs("args", args)...)
	if _, err := p.ProcessDetailsUpdater.UpdateProcessDetails(ctx, &repo.UpdateProcessDetailsRequest{
		ID:      details.ID,
		Command: ptr.To(strings.Join(args, " ")),
	}); err != nil {
		_ = logFile.Close()
		return withBaseErr(err, "failed to update process details(%d) command", details.ID)
	}

	var (
		resultErrors = []error{}
		addErr       = func(err error) {
			if err == nil {
				return
			}
			resultErrors = append(resultErrors, err)
		}
	)

	cmd := exec.CommandContext(ctx, p.Pneutrinoutil, args...)
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	cmd.Env = p.Env
	err = cmd.Run()
	_ = logFile.Close()
	if err != nil {
		addErr(withBaseErr(err, "failed to run pneutrinoutil process"))
	}

	resultObjectPath := filepath.Join(p.BasePath, payload.RequestID)
	alog.L().Info("upload log", attrs("from", logPath, "to", resultObjectPath)...)

	logObjectId := func() *int {
		logObjectId, err := p.uploadLog(ctx, logPath, resultObjectPath)
		if err != nil {
			addErr(withBaseErr(err, "failed to upload log"))
			return nil
		}
		return &logObjectId
	}()

	resultObjectId := func() *int {
		alog.L().Info("find local results directory", attrs("from", workDir)...)
		resultDir, err := p.findResultDir(workDir)
		if err != nil {
			addErr(withBaseErr(err, "failed to find local results directory"))
			return nil
		}
		alog.L().Info("upload results", attrs("from", resultDir, "to", resultObjectPath)...)
		resultObjectId, err := p.uploadResults(ctx, resultDir, resultObjectPath, filepath.Base(score.Path))
		if err != nil {
			addErr(withBaseErr(err, "failed to upload results"))
			return nil
		}
		return &resultObjectId
	}()

	alog.L().Info("update log and results in process details", attrs("id", details.ID)...)
	if _, err := p.ProcessDetailsUpdater.UpdateProcessDetails(ctx, &repo.UpdateProcessDetailsRequest{
		ID:             details.ID,
		LogObjectId:    logObjectId,
		ResultObjectId: resultObjectId,
	}); err != nil {
		addErr(withBaseErr(err, "failed to update process_details(%d)", details.ID))
	}

	if len(resultErrors) > 0 {
		return errors.Join(resultErrors...)
	}

	alog.L().Info("succeed", attrs()...)
	processSucceed = true
	return nil
}

func (p *PneutrinoutilProcessor) generateArgs(payload *PneutrinoutilStartPayload, workDir, scorePath string) []string {
	args := append([]string{
		"--desc", payload.RequestID,
		"--neutrinoDir", p.NeutrinoDir,
		"--workDir", workDir,
		"--score", scorePath,
		"--env", "all",
		"--shell", p.Shell,
	}, payload.Args...)
	escaped := make([]string, len(args))
	for i, x := range args {
		escaped[i] = shellescape.Quote(x)
	}
	return escaped
}

func (p *PneutrinoutilProcessor) uploadLog(ctx context.Context, logPath, resultObjectPath string) (int, error) {
	f, err := os.Open(logPath)
	if err != nil {
		return 0, err
	}
	defer func() { _ = f.Close() }()

	path := filepath.Join(resultObjectPath, filepath.Base(logPath))
	r, err := p.writeObject(ctx, path, f)
	if err != nil {
		return 0, err
	}
	return r.ID, nil
}

func (p *PneutrinoutilProcessor) uploadResults(ctx context.Context, resultDir, resultObjectPath, ignore string) (int, error) {
	elems, err := os.ReadDir(resultDir)
	if err != nil {
		return 0, err
	}
	for _, elem := range elems {
		if elem.IsDir() {
			continue
		}
		if elem.Name() == ignore {
			continue
		}
		path := filepath.Join(resultDir, elem.Name())
		f, err := os.Open(path)
		if err != nil {
			return 0, err
		}
		_, err = p.writeObject(ctx, filepath.Join(resultObjectPath, filepath.Base(path)), f)
		_ = f.Close()
		if err != nil {
			return 0, err
		}
	}

	r, err := p.writeObject(ctx, resultObjectPath, nil)
	if err != nil {
		return 0, err
	}
	return r.ID, nil
}

func (p *PneutrinoutilProcessor) writeObject(ctx context.Context, path string, blob io.ReadSeeker) (*domain.Object, error) {
	typ := domain.ObjectTypeDir
	if blob != nil {
		typ = domain.ObjectTypeFile
	}
	r, err := p.ObjectWriter.WriteObject(ctx, &repo.WriteObjectRequest{
		Type:   typ,
		Bucket: p.Bucket,
		Path:   path,
		Blob:   blob,
	})
	logAttrs := []any{
		"task", TypePneutrinoutilStart,
		"bucket", p.Bucket,
		"path", path,
	}

	if err != nil {
		alog.L().Debug("write object", append(logAttrs, logx.Err(err))...)
		return nil, err
	}
	alog.L().Debug("write object", append(logAttrs, "sizeBytes", r.Object().SizeBytes, "id", r.Object().ID)...)
	return r.Object(), nil
}

func (p *PneutrinoutilProcessor) findResultDir(workDir string) (string, error) {
	dirs, err := os.ReadDir(filepath.Join(workDir, "result"))
	if err != nil {
		return "", fmt.Errorf("%w: failed to read result dir", err)
	}
	if len(dirs) != 1 {
		return "", fmt.Errorf("failed to read result dir: want single dir")
	}
	d := dirs[0]
	if !d.IsDir() {
		return "", fmt.Errorf("failed to read result internal dir: want single dir")
	}
	return filepath.Join(workDir, "result", d.Name()), nil
}

func (p *PneutrinoutilProcessor) updateProcessStatus(ctx context.Context, processID int, now time.Time, processSucceed bool) error {
	status := ptr.To(domain.ProcessStatusFailed)
	if processSucceed {
		status = ptr.To(domain.ProcessStatusSucceed)
	}
	_, err := p.ProcessUpdater.UpdateProcess(ctx, &repo.UpdateProcessRequest{
		ID:          processID,
		Status:      status,
		CompletedAt: ptr.To(now),
	})
	return err
}

func (p *PneutrinoutilProcessor) webhook(ctx context.Context, requestID string, processSucceed bool) error {
	if p.Webhooker == nil {
		return nil
	}
	return p.Webhooker.Webhook(ctx, &PneutrinoutilStartWebhookParams{
		Type:      TypePneutrinoutilStart,
		RequestID: requestID,
		OK:        processSucceed,
	})
}
