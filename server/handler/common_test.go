package handler_test

import (
	"testing"
	"time"

	"github.com/berquerant/pneutrinoutil/server/handler"
	"github.com/stretchr/testify/assert"
)

func TestCustomTime(t *testing.T) {
	t.Run("UnmarshalParam", func(t *testing.T) {
		for _, tc := range []struct {
			title string
			param string
			want  handler.CustomTime
		}{
			{
				title: "timestamp",
				param: "1763188414",
				want:  handler.CustomTime(time.Unix(1763188414, 0)),
			},
			{
				title: "rfc3339",
				param: "2025-11-15T15:33:34+09:00",
				want:  handler.CustomTime(time.Unix(1763188414, 0)),
			},
		} {
			t.Run(tc.title, func(t *testing.T) {
				got := new(handler.CustomTime)
				err := got.UnmarshalParam(tc.param)
				assert.Nil(t, err, "%v", err)
				assert.Equal(t, tc.want.Timestamp(), got.Timestamp())
			})
		}
	})
}
