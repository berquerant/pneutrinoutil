// Code generated by "dataclass -type BaseResult -field RequestID string -output result_baseresult_dataclass_generated.go"; DO NOT EDIT.

package pworker

type BaseResult interface {
	RequestID() string
}
type baseResult struct {
	requestID string
}

func (s *baseResult) RequestID() string { return s.requestID }
func NewBaseResult(
	requestID string,
) BaseResult {
	return &baseResult{
		requestID: requestID,
	}
}
