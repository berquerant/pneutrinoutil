package pathx_test

import (
	"testing"
	"time"

	"github.com/berquerant/pneutrinoutil/pkg/pathx"
	"github.com/stretchr/testify/assert"
)

func TestResultElement(t *testing.T) {
	for _, tc := range []struct {
		e *pathx.ResultElement
		s string
	}{
		{
			e: pathx.NewResultElement(
				"BASE",
				time.Date(2009, time.November, 10, 23, 1, 2, 0, time.UTC),
				int64(1257894062),
				100,
			),
			s: "BASE__20091110230102_1257894062_100",
		},
	} {
		t.Run(tc.e.String(), func(t *testing.T) {
			assert.Equal(t, tc.s, tc.e.String())
			got, err := pathx.ParseResultElement(tc.s)
			assert.Nil(t, err)
			assert.Equal(t, tc.e, got)
		})
	}
}
