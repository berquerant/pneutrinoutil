package repo

type Range[T any] struct {
	Left  *T
	Right *T
}

func NewRange[T any](left, right *T) Range[T] {
	return Range[T]{
		Left:  left,
		Right: right,
	}
}
