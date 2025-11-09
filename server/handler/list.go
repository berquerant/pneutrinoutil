package handler

import (
	"net/http"
	"time"

	"github.com/berquerant/pneutrinoutil/pkg/domain"
	"github.com/berquerant/pneutrinoutil/pkg/repo"
	"github.com/labstack/echo/v4"
)

type List struct {
	processLister repo.ProcessLister
	detailsGetter repo.ProcessDetailsGetter
}

func NewList(processLister repo.ProcessLister, detailsGetter repo.ProcessDetailsGetter) *List {
	return &List{
		processLister: processLister,
		detailsGetter: detailsGetter,
	}
}

type ListProcessParam struct {
	Limit  int    `query:"limit"`  // default: 5
	Status string `query:"status"` // (pending|running|succeed|failed)
}

type ListProcessResponseData []*ListProcessResponseDataElement

func (d ListProcessResponseData) Len() int { return len(d) }

type ListProcessResponseDataElement struct {
	RequestID      string `json:"rid"` // request id, or just id
	DetailsID      int    `json:"details_id"`
	Basename       string `json:"basename"` // title
	Command        string `json:"command,omitempty"`
	Status         string `json:"status"`
	CreatedAt      string `json:"created_at"`
	StartedAt      string `json:"started_at,omitempty"`
	CompletedAt    string `json:"completed_at,omitempty"`
	ScoreObjectID  int    `json:"score_object_id"`
	LogObjectID    int    `json:"log_object_id,omitempty"`
	ResultObjectID int    `json:"result_object_id,omitempty"`
}

func (p ListProcessParam) intoRequest() *repo.ListProcessRequest {
	var r repo.ListProcessRequest
	if x := p.Limit; x > 0 {
		r.Limit = x
	} else {
		r.Limit = 5
	}
	if x, ok := domain.ProcessStatusFromString(p.Status); ok {
		r.Status = &x
	}
	return &r
}

// List process results
//
// @summary list results
// @description list results of processes
// @param limit query int false "query limit; default: 5"
// @param status query string false "process status; (pending|running|succeed|failed)"
// @produce json
// @success 200 {object} handler.SuccessResponse[ListProcessResponseData]
// @router /proc [get]
func (s *List) ListProcess(c echo.Context) error {
	var p ListProcessParam
	if err := c.Bind(&p); err != nil {
		return err
	}

	xs, err := s.processLister.ListProcess(c.Request().Context(), p.intoRequest())
	if err != nil {
		return err
	}

	ids := make([]int, len(xs))                                     // details_id
	idMap := make(map[int]*ListProcessResponseDataElement, len(xs)) // details_id -> element
	data := make([]*ListProcessResponseDataElement, len(xs))
	for i, x := range xs {
		ids[i] = x.DetailsID

		y := &ListProcessResponseDataElement{
			RequestID: x.RequestID,
			DetailsID: x.DetailsID,
			Status:    x.Status.String(),
			CreatedAt: x.CreatedAt.Format(time.DateTime),
		}
		if t := x.StartedAt; t != nil {
			y.StartedAt = t.Format(time.DateTime)
		}
		if t := x.CompletedAt; t != nil {
			y.CompletedAt = t.Format(time.DateTime)
		}

		data[i] = y
		idMap[x.DetailsID] = y
	}

	ys, err := s.detailsGetter.GetProcessDetailsList(c.Request().Context(), ids...)
	if err != nil {
		return err
	}
	for _, x := range ys {
		if y, ok := idMap[x.ID]; ok {
			y.Basename = x.Title
			if v := x.Command; v != nil {
				y.Command = *v
			}
			y.ScoreObjectID = x.ScoreObjectID
			if v := x.LogObjectID; v != nil {
				y.LogObjectID = *v
			}
			if v := x.ResultObjectID; v != nil {
				y.ResultObjectID = *v
			}
		}
	}

	return Success(c, http.StatusOK, ListProcessResponseData(data))
}
