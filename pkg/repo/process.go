package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/berquerant/pneutrinoutil/pkg/domain"
	"github.com/berquerant/pneutrinoutil/pkg/infra"
	"github.com/berquerant/pneutrinoutil/pkg/ptr"
)

const (
	ProcessTable = "processes"
)

type CreateProcessRequest struct {
	RequestId   string
	Status      domain.ProcessStatus
	DetailsId   int
	StartedAt   *time.Time
	CompletedAt *time.Time
}

type ProcessCreator interface {
	CreateProcess(ctx context.Context, req *CreateProcessRequest) (*domain.Process, error)
}

type UpdateProcessRequest struct {
	ID          int
	Status      *domain.ProcessStatus
	StartedAt   *time.Time
	CompletedAt *time.Time
}

type ProcessUpdater interface {
	UpdateProcess(ctx context.Context, req *UpdateProcessRequest) (*domain.Process, error)
}

type ProcessGetter interface {
	GetProcess(ctx context.Context, id int) (*domain.Process, error)
	GetProcessByRequestId(ctx context.Context, rid string) (*domain.Process, error)
	GetProcessByDetailsList(ctx context.Context, detailsID ...int) ([]*domain.Process, error)
}

type ListProcessRequest struct {
	Limit  int
	Status *domain.ProcessStatus
}

type ProcessLister interface {
	ListProcess(ctx context.Context, req *ListProcessRequest) ([]*domain.Process, error)
}

var (
	_ ProcessCreator = &Process{}
	_ ProcessUpdater = &Process{}
	_ ProcessGetter  = &Process{}
	_ ProcessLister  = &Process{}
)

func NewProcess(query infra.Queryer[domain.Process], exec infra.Execer) *Process {
	return &Process{
		query: query,
		exec:  exec,
	}
}

type Process struct {
	query infra.Queryer[domain.Process]
	exec  infra.Execer
}

func (p *Process) UpdateProcess(ctx context.Context, req *UpdateProcessRequest) (*domain.Process, error) {
	var (
		args []any
		cols []string
	)

	if x := req.Status; x != nil {
		cols = append(cols, "status_id = ?")
		args = append(args, int(*x))
	}
	if x := req.StartedAt; x != nil {
		cols = append(cols, "started_at = ?")
		args = append(args, *x)
	}
	if x := req.CompletedAt; x != nil {
		cols = append(cols, "completed_at = ?")
		args = append(args, *x)
	}

	if _, err := p.exec.Exec(ctx, &infra.ExecRequest{
		Query: fmt.Sprintf("update processes set %s where id = ?;", strings.Join(cols, ",")),
		Args:  append(args, req.ID),
	}); err != nil {
		return nil, fmt.Errorf("%w: update process: id=%d", err, req.ID)
	}

	return p.GetProcess(ctx, req.ID)
}

func (p *Process) CreateProcess(ctx context.Context, req *CreateProcessRequest) (*domain.Process, error) {
	r, err := p.exec.Exec(ctx, &infra.ExecRequest{
		Query: "insert into processes (request_id, status_id, details_id, started_at, completed_at) values (?, ?, ?, ?, ?);",
		Args: []any{
			req.RequestId,
			int(req.Status),
			req.DetailsId,
			req.StartedAt,
			req.CompletedAt,
		},
		AssertResponse: infra.AssertJoin(
			infra.AssertRowsAffected(1),
			infra.AssertLastInserted(),
		),
	})
	if err != nil {
		return nil, fmt.Errorf("%w: create process: request_id=%s", err, req.RequestId)
	}

	id, _ := r.LastInsertId()
	return p.GetProcess(ctx, int(id))
}

func (*Process) scan(f func(...any) error) (*domain.Process, error) {
	var (
		id          int
		requestId   string
		statusId    int
		detailsId   int
		startedAt   sql.NullTime
		completedAt sql.NullTime
		createdAt   time.Time
		updatedAt   time.Time
	)
	if err := f(&id, &requestId, &statusId, &detailsId, &startedAt, &completedAt, &createdAt, &updatedAt); err != nil {
		return nil, err
	}
	v := &domain.Process{
		ID:        id,
		RequestID: requestId,
		Status:    domain.ProcessStatus(statusId),
		DetailsID: detailsId,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
	if startedAt.Valid {
		v.StartedAt = ptr.To(startedAt.Time)
	}
	if completedAt.Valid {
		v.CompletedAt = ptr.To(completedAt.Time)
	}
	return v, nil
}

func (p *Process) GetProcess(ctx context.Context, id int) (*domain.Process, error) {
	r, err := p.query.Query(ctx, &infra.QueryRequest[domain.Process]{
		Query: "select id, request_id, status_id, details_id, started_at, completed_at, created_at, updated_at from processes where id = ?;",
		Args: []any{
			id,
		},
		Scan: p.scan,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: get process: id=%d", err, id)
	}
	if err := r.AssertRows(1); err != nil {
		return nil, fmt.Errorf("%w: get process: id=%d", err, id)
	}
	return r.Items[0], nil
}

func (p *Process) GetProcessByRequestId(ctx context.Context, rid string) (*domain.Process, error) {
	r, err := p.query.Query(ctx, &infra.QueryRequest[domain.Process]{
		Query: "select id, request_id, status_id, details_id, started_at, completed_at, created_at, updated_at from processes where request_id = ?;",
		Args: []any{
			rid,
		},
		Scan: p.scan,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: get process by request id %s", err, rid)
	}
	if err := r.AssertRows(1); err != nil {
		return nil, fmt.Errorf("%w: get process by request id %s", err, rid)
	}
	return r.Items[0], nil
}

func (p *Process) GetProcessByDetailsList(ctx context.Context, detailsID ...int) ([]*domain.Process, error) {
	if len(detailsID) == 0 {
		return nil, nil
	}

	xs := make([]string, len(detailsID))
	for i, v := range detailsID {
		xs[i] = fmt.Sprint(v)
	}
	r, err := p.query.Query(ctx, &infra.QueryRequest[domain.Process]{
		Query: fmt.Sprintf("select id, request_id, status_id, details_id, started_at, completed_at, created_at, updated_at from processes where details_id in (%s);",
			strings.Join(xs, ","),
		),
		Scan: p.scan,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: get process list by details_id: id=%v", err, detailsID)
	}
	return r.Items, nil
}

func (p *Process) ListProcess(ctx context.Context, req *ListProcessRequest) ([]*domain.Process, error) {
	var (
		where string
		args  []any
	)
	if x := req.Status; x != nil {
		where = "where status = ?"
		args = append(args, int(*x))
	}
	if x := req.Limit; x > 0 {
		args = append(args, x)
	} else {
		args = append(args, 5)
	}
	r, err := p.query.Query(ctx, &infra.QueryRequest[domain.Process]{
		Query: fmt.Sprintf("select id, request_id, status_id, details_id, started_at, completed_at, created_at, updated_at from processes %s order by id desc limit ?;", where),
		Args:  args,
		Scan:  p.scan,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: list process", err)
	}
	return r.Items, nil
}
