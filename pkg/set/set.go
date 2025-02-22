package set

import "maps"

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

func (s Set[T]) Diff(other Set[T]) Set[T] {
	r := maps.Clone(s)
	maps.DeleteFunc(r, func(key T, value bool) bool {
		if !value {
			// delete garbages
			return true
		}
		return other.In(key)
	})
	return r
}

func (s Set[T]) IntoSlice() []T {
	r := []T{}
	for key, ok := range s {
		if ok {
			r = append(r, key)
		}
	}
	return r
}
