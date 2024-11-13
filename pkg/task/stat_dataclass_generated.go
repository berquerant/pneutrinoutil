// Code generated by "dataclass -type Stat -field Title string|Elapsed time.Duration -output stat_dataclass_generated.go"; DO NOT EDIT.

package task

import "time"

type Stat interface {
	Title() string
	Elapsed() time.Duration
}
type stat struct {
	title   string
	elapsed time.Duration
}

func (s *stat) Title() string          { return s.title }
func (s *stat) Elapsed() time.Duration { return s.elapsed }
func NewStat(
	title string,
	elapsed time.Duration,
) Stat {
	return &stat{
		title:   title,
		elapsed: elapsed,
	}
}