package functional

// Filter implements slice filtering functionality
func Filter[T any](s []T, f func(T) bool) []T {
	m := make([]T, 0)
	for _, e := range s {
		if f(e) {
			m = append(m, e)
		}
	}
	return m
}
