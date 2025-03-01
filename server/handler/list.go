package handler

import (
	"net/http"
	"time"

	"github.com/berquerant/pneutrinoutil/server/pworker"
	"github.com/labstack/echo/v4"
)

type List struct {
	list *pworker.List
}

func NewList(list *pworker.List) *List {
	return &List{
		list: list,
	}
}

// List process results
//
// @summary list results
// @description list results of processes
// @produce json
// @success 200 {object} handler.SuccessResponse[ListResponseData]
// @router /proc [get]
func (s *List) Handler(c echo.Context) error {
	keys := s.list.ListKeys()
	data := make([]*ListResponseDataElement, len(keys))
	for i, x := range keys {
		data[i] = &ListResponseDataElement{
			RequestID: x.RequestID,
			Basename:  x.Element.Basename,
			CreatedAt: x.Element.Now.Format(time.DateTime),
			Salt:      x.Element.Salt,
		}
	}
	return Success(c, http.StatusOK, ListResponseData(data))
}

type ListResponseData []*ListResponseDataElement

func (d ListResponseData) Len() int { return len(d) }

type ListResponseDataElement struct {
	RequestID string `json:"rid"`      // request id, or just id
	Basename  string `json:"basename"` // original musicxml file name except extension
	CreatedAt string `json:"created_at"`
	Salt      int    `json:"salt"`
}
