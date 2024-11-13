package set

func New[T comparable](xs []T) Set[T] {
	d := map[T]bool{}
	for _, x := range xs {
		d[x] = true
	}
	return Set[T](d)
}

type Set[T comparable] map[T]bool

func (s Set[T]) In(t T) bool { return s[t] }
func (s Set[T]) Len() int    { return len(s) }
