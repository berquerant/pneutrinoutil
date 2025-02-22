package set_test

import (
	"slices"
	"testing"

	"github.com/berquerant/pneutrinoutil/pkg/set"
	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	t.Run("Diff", func(t *testing.T) {
		for _, tc := range []struct {
			title       string
			left, right set.Set[int]
			want        []int
		}{
			{
				title: "empty-empty",
				left:  set.New([]int{}),
				right: set.New([]int{}),
				want:  []int{},
			},
			{
				title: "empty-some",
				left:  set.New([]int{}),
				right: set.New([]int{1}),
				want:  []int{},
			},
			{
				title: "some-empty",
				left:  set.New([]int{1}),
				right: set.New([]int{}),
				want:  []int{1},
			},
			{
				title: "some-some",
				left:  set.New([]int{1, 2, 3}),
				right: set.New([]int{2, 4}),
				want:  []int{1, 3},
			},
		} {
			t.Run(tc.title, func(t *testing.T) {
				g := tc.left.Diff(tc.right).IntoSlice()
				slices.Sort(g)
				slices.Sort(tc.want)
				assert.Equal(t, tc.want, g)
			})
		}
	})
}
