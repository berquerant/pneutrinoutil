package ptr

//go:fix inline
func To[T any](v T) *T {
	return new(v)
}
