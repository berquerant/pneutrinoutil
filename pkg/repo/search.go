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

type SearchProcessRequest struct {
	Limit       int
	Status      *domain.ProcessStatus
	TitlePrefix *string
	CreatedAt   Range[time.Time]
}

type SearchProcessResult struct {
	Items []*SearchProcessResultElement
}

type SearchProcessResultElement struct {
	Process *domain.Process
	Details *domain.ProcessDetails
}

type ProcessSearcher interface {
	SearchProcess(ctx context.Context, req *SearchProcessRequest) (*SearchProcessResult, error)
}

var _ ProcessSearcher = &Searcher{}

func NewSearcher(query infra.Queryer[SearchProcessResultElement]) *Searcher {
	return &Searcher{
		query: query,
	}
}

type Searcher struct {
	query infra.Queryer[SearchProcessResultElement]
}

func (*Searcher) scanSearchProcessResultElement(f func(...any) error) (*SearchProcessResultElement, error) {
	var (
		processId        int
		requestId        string
		statusId         int
		detailsId        int
		startedAt        sql.NullTime
		completedAt      sql.NullTime
		processCreatedAt time.Time
		processUpdatedAt time.Time
		__detailsId      int
		command          sql.NullString
		title            string
		scoreObjectId    int
		logObjectId      sql.NullInt64
		resultObjectId   sql.NullInt64
		detailsCreatedAt time.Time
		detailsUpdatedAt time.Time
	)
	if err := f(
		&processId, &requestId, &statusId, &detailsId, &startedAt, &completedAt, &processCreatedAt, &processUpdatedAt,
		&__detailsId, &command, &title, &scoreObjectId, &logObjectId, &resultObjectId, &detailsCreatedAt, &detailsUpdatedAt,
	); err != nil {
		return nil, err
	}

	p := &domain.Process{
		ID:        processId,
		RequestID: requestId,
		Status:    domain.ProcessStatus(statusId),
		DetailsID: detailsId,
		CreatedAt: processCreatedAt,
		UpdatedAt: processUpdatedAt,
	}
	if startedAt.Valid {
		p.StartedAt = ptr.To(startedAt.Time)
	}
	if completedAt.Valid {
		p.CompletedAt = ptr.To(completedAt.Time)
	}

	d := &domain.ProcessDetails{
		ID:            detailsId,
		Title:         title,
		ScoreObjectID: scoreObjectId,
		CreatedAt:     detailsCreatedAt,
		UpdatedAt:     detailsUpdatedAt,
	}
	if command.Valid {
		d.Command = ptr.To(command.String)
	}
	if logObjectId.Valid {
		d.LogObjectID = ptr.To(int(logObjectId.Int64))
	}
	if resultObjectId.Valid {
		d.ResultObjectID = ptr.To(int(resultObjectId.Int64))
	}

	return &SearchProcessResultElement{
		Process: p,
		Details: d,
	}, nil
}

func (s *Searcher) SearchProcess(ctx context.Context, req *SearchProcessRequest) (*SearchProcessResult, error) {
	const baseQuery = `select
p.id, p.request_id, p.status_id, p.details_id, p.started_at, p.completed_at, p.created_at, p.updated_at,
d.id, d.command, d.title, d.score_object_id, d.log_object_id, d.result_object_id, d.created_at, d.updated_at
from process_details d inner join processes p on d.id = p.details_id`
	var (
		conditions []string
		args       []any
	)
	if x := req.Status; x != nil {
		conditions = append(conditions, "p.status_id = ?")
		args = append(args, *x)
	}
	if x := req.TitlePrefix; x != nil {
		conditions = append(conditions, "d.title like ?")
		args = append(args, *x+"%")
	}
	if x := req.CreatedAt.Left; x != nil {
		conditions = append(conditions, "p.created_at > ?")
		args = append(args, *x)
	}
	if x := req.CreatedAt.Right; x != nil {
		conditions = append(conditions, "p.created_at < ?")
		args = append(args, *x)
	}

	query := baseQuery
	if len(conditions) > 0 {
		query += " where " + strings.Join(conditions, " and ")
	}
	query += " order by p.created_at desc limit ?;"
	args = append(args, req.Limit)

	r, err := s.query.Query(ctx, &infra.QueryRequest[SearchProcessResultElement]{
		Query: query,
		Args:  args,
		Scan:  s.scanSearchProcessResultElement,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: search process", err)
	}
	return &SearchProcessResult{
		Items: r.Items,
	}, nil
}
