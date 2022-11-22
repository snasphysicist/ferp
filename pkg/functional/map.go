package functional

// Map implements slice mapping functionality
func Map[T any, U any](s []T, f func(T) U) []U {
	m := make([]U, len(s))
	for i, e := range s {
		m[i] = f(e)
	}
	return m
}
