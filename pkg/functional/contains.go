package functional

// Contains returns true iff the slice contains the element
func Contains[T comparable](a []T, e T) bool {
	return len(Filter(a, func(ae T) bool { return ae == e })) > 0
}
