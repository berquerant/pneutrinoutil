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
	ProcessDetailsTable = "process_details"
)

type CreateProcessDetailsRequest struct {
	Command        *string
	Title          string
	ScoreObjectId  int
	LogObjectId    *int
	ResultObjectId *int
}

type ProcessDetailsCreator interface {
	CreateProcessDetails(ctx context.Context, req *CreateProcessDetailsRequest) (*domain.ProcessDetails, error)
}

type UpdateProcessDetailsRequest struct {
	ID             int
	Command        *string
	LogObjectId    *int
	ResultObjectId *int
}

type ProcessDetailsUpdater interface {
	UpdateProcessDetails(ctx context.Context, req *UpdateProcessDetailsRequest) (*domain.ProcessDetails, error)
}

type ProcessDetailsGetter interface {
	GetProcessDetails(ctx context.Context, id int) (*domain.ProcessDetails, error)
	GetProcessDetailsList(ctx context.Context, id ...int) ([]*domain.ProcessDetails, error)
}

var (
	_ ProcessDetailsCreator = &ProcessDetails{}
	_ ProcessDetailsUpdater = &ProcessDetails{}
	_ ProcessDetailsGetter  = &ProcessDetails{}
)

func NewProcessDetails(query infra.Queryer[domain.ProcessDetails], exec infra.Execer) *ProcessDetails {
	return &ProcessDetails{
		query: query,
		exec:  exec,
	}
}

type ProcessDetails struct {
	query infra.Queryer[domain.ProcessDetails]
	exec  infra.Execer
}

func (p *ProcessDetails) UpdateProcessDetails(ctx context.Context, req *UpdateProcessDetailsRequest) (*domain.ProcessDetails, error) {
	var (
		args []any
		cols []string
	)
	if x := req.Command; x != nil {
		cols = append(cols, "command = ?")
		args = append(args, *x)
	}
	if x := req.LogObjectId; x != nil {
		cols = append(cols, "log_object_id = ?")
		args = append(args, *x)
	}
	if x := req.ResultObjectId; x != nil {
		cols = append(cols, "result_object_id = ?")
		args = append(args, *x)
	}

	if _, err := p.exec.Exec(ctx, &infra.ExecRequest{
		Query: fmt.Sprintf("update process_details set %s where id = ?;", strings.Join(cols, ",")),
		Args:  append(args, req.ID),
	}); err != nil {
		return nil, fmt.Errorf("%w: update process details: id=%d", err, req.ID)
	}

	return p.GetProcessDetails(ctx, req.ID)
}

func (p *ProcessDetails) CreateProcessDetails(ctx context.Context, req *CreateProcessDetailsRequest) (*domain.ProcessDetails, error) {
	r, err := p.exec.Exec(ctx, &infra.ExecRequest{
		Query: "insert into process_details (command, title, score_object_id, log_object_id, result_object_id) values (?, ?, ?, ?, ?);",
		Args: []any{
			req.Command,
			req.Title,
			req.ScoreObjectId,
			req.LogObjectId,
			req.ResultObjectId,
		},
		AssertResponse: infra.AssertJoin(
			infra.AssertRowsAffected(1),
			infra.AssertLastInserted(),
		),
	})
	if err != nil {
		return nil, fmt.Errorf("%w: create process details: title=%s", err, req.Title)
	}

	id, _ := r.LastInsertId()
	return p.GetProcessDetails(ctx, int(id))
}

func (*ProcessDetails) scan(f func(...any) error) (*domain.ProcessDetails, error) {
	var (
		id             int
		command        sql.NullString
		title          string
		scoreObjectId  int
		logObjectId    sql.NullInt64
		resultObjectId sql.NullInt64
		createdAt      time.Time
		updatedAt      time.Time
	)
	if err := f(&id, &command, &title, &scoreObjectId, &logObjectId, &resultObjectId, &createdAt, &updatedAt); err != nil {
		return nil, err
	}
	v := &domain.ProcessDetails{
		ID:            id,
		Title:         title,
		ScoreObjectID: scoreObjectId,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
	if command.Valid {
		v.Command = ptr.To(command.String)
	}
	if logObjectId.Valid {
		v.LogObjectID = ptr.To(int(logObjectId.Int64))
	}
	if resultObjectId.Valid {
		v.ResultObjectID = ptr.To(int(resultObjectId.Int64))
	}
	return v, nil
}

func (p *ProcessDetails) GetProcessDetails(ctx context.Context, id int) (*domain.ProcessDetails, error) {
	r, err := p.query.Query(ctx, &infra.QueryRequest[domain.ProcessDetails]{
		Query: "select id, command, title, score_object_id, log_object_id, result_object_id, created_at, updated_at from process_details where id = ?",
		Args: []any{
			id,
		},
		Scan: p.scan,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: get process details: id=%d", err, id)
	}
	if err := r.AssertRows(1); err != nil {
		return nil, fmt.Errorf("%w: get process details: ud=%d", err, id)
	}
	return r.Items[0], nil
}

func (p *ProcessDetails) GetProcessDetailsList(ctx context.Context, id ...int) ([]*domain.ProcessDetails, error) {
	if len(id) == 0 {
		return nil, nil
	}

	xs := make([]string, len(id))
	for i, v := range id {
		xs[i] = fmt.Sprint(v)
	}
	r, err := p.query.Query(ctx, &infra.QueryRequest[domain.ProcessDetails]{
		Query: fmt.Sprintf("select id, command, title, score_object_id, log_object_id, result_object_id, created_at, updated_at from process_details where id in (%s);",
			strings.Join(xs, ","),
		),
		Scan: p.scan,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: get process details list: id=%v", err, id)
	}
	return r.Items, nil
}
