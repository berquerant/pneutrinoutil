package handler

import (
	"net/http"
	"time"

	"github.com/berquerant/pneutrinoutil/pkg/domain"
	"github.com/berquerant/pneutrinoutil/pkg/ptr"
	"github.com/berquerant/pneutrinoutil/pkg/repo"
	"github.com/labstack/echo/v4"
)

type Search struct {
	processSearcher repo.ProcessSearcher
}

func NewSearch(processSearcher repo.ProcessSearcher) *Search {
	return &Search{
		processSearcher: processSearcher,
	}
}

type SearchProcessParam struct {
	Limit  int         `query:"limit"`  // default: 5
	Status string      `query:"status"` // (pending|running|succeed|failed)
	Prefix string      `query:"prefix"` // title prefix
	Start  *CustomTime `query:"start"`  // created_at; RFC3339 or timestamp
	End    *CustomTime `query:"end"`    // created_at; RFC3339 or timestamp
}

func (p SearchProcessParam) intoRequest() *repo.SearchProcessRequest {
	var r repo.SearchProcessRequest
	if x := p.Limit; x > 0 {
		r.Limit = x
	} else {
		r.Limit = 5
	}
	if x, ok := domain.ProcessStatusFromString(p.Status); ok {
		r.Status = &x
	}
	if x := p.Prefix; x != "" {
		r.TitlePrefix = &x
	}
	var (
		start, end *time.Time
	)
	if x := p.Start; x != nil {
		start = ptr.To(time.Time(*x))
	}
	if x := p.End; x != nil {
		end = ptr.To(time.Time(*x))
	}
	r.CreatedAt = repo.NewRange(start, end)
	return &r
}

type SearchProcessResponseDataElement struct {
	RequestID      string    `json:"request_id"` // request id, or just id
	Status         string    `json:"status"`
	DetailsID      int       `json:"details_id"`
	StartedAt      time.Time `json:"started_at,omitempty"`
	CompletedAt    time.Time `json:"completed_at,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Command        string    `json:"command,omitempty"`
	Title          string    `json:"title"`
	ScoreObjectID  int       `json:"score_object_id"`
	LogObjectID    int       `json:"log_object_id,omitempty"`
	ResultObjectID int       `json:"result_object_id,omitempty"`
}

type SearchProcessResponseData []*SearchProcessResponseDataElement

// Search processes.
//
// @sumary search processes
// @description search processes by status, created_at, title prefix, order by created_at
// @param limit query int false "query limit; default: 5"
// @param prefix query string false "title prefix"
// @param status query string false "process status; (pending|running|succeed|failed)"
// @param start query string false "created_at"
// @param end query string false "created_at"
// @produce json
// @success 200 {object} handler.SuccessResponse[SearchProcessResponseData]
// @router /proc/search [get]
func (s *Search) SearchProcess(c echo.Context) error {
	var p SearchProcessParam
	if err := c.Bind(&p); err != nil {
		return err
	}

	xs, err := s.processSearcher.SearchProcess(c.Request().Context(), p.intoRequest())
	if err != nil {
		return err
	}

	data := make([]*SearchProcessResponseDataElement, len(xs.Items))
	for i, x := range xs.Items {
		y := &SearchProcessResponseDataElement{
			RequestID:     x.Process.RequestID,
			Status:        x.Process.Status.String(),
			DetailsID:     x.Process.DetailsID,
			CreatedAt:     x.Process.CreatedAt,
			UpdatedAt:     x.Process.UpdatedAt,
			Title:         x.Details.Title,
			ScoreObjectID: x.Details.ScoreObjectID,
		}
		if v := x.Process.StartedAt; v != nil {
			y.StartedAt = *v
		}
		if v := x.Process.CompletedAt; v != nil {
			y.CompletedAt = *v
		}
		if v := x.Details.Command; v != nil {
			y.Command = *v
		}
		if v := x.Details.LogObjectID; v != nil {
			y.LogObjectID = *v
		}
		if v := x.Details.ResultObjectID; v != nil {
			y.ResultObjectID = *v
		}
		data[i] = y
	}

	return Success(c, http.StatusOK, SearchProcessResponseData(data))
}
