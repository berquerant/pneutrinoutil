package pworker

import (
	"github.com/berquerant/pneutrinoutil/cli/ctl"
	"github.com/berquerant/pneutrinoutil/pkg/pathx"
)

type Result interface {
	RequestID() string
	Err() error
	Dir() string
	Config() *ctl.Config
	Element() *pathx.ResultElement
}

var (
	_ Result = &ErrorResult{}
	_ Result = &SuccessResult{}
)

//go:generate go tool dataclass -type "BaseResult" -field "RequestID string" -output result_baseresult_dataclass_generated.go

func NewErrorResult(rid string, err error) *ErrorResult {
	return &ErrorResult{
		BaseResult: NewBaseResult(rid),
		err:        err,
	}
}

type ErrorResult struct {
	BaseResult
	err error
}

func (r *ErrorResult) Err() error                  { return r.err }
func (*ErrorResult) Dir() string                   { return "" }
func (*ErrorResult) Config() *ctl.Config           { return nil }
func (*ErrorResult) Element() *pathx.ResultElement { return nil }

func NewSuccessResult(rid, dir string, c *ctl.Config, elem *pathx.ResultElement) *SuccessResult {
	return &SuccessResult{
		BaseResult: NewBaseResult(rid),
		dir:        dir,
		config:     c,
		element:    elem,
	}
}

type SuccessResult struct {
	BaseResult
	dir     string
	config  *ctl.Config
	element *pathx.ResultElement
}

func (*SuccessResult) Err() error                      { return nil }
func (r *SuccessResult) Dir() string                   { return r.dir }
func (r *SuccessResult) Config() *ctl.Config           { return r.config }
func (r *SuccessResult) Element() *pathx.ResultElement { return r.element }
