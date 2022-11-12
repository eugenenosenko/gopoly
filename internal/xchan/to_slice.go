package xchan

func ToSlice[O ~[]T, T any](in <-chan T) O {
	out := make([]T, 0, len(in))
	for i := range in {
		out = append(out, i)
	}
	return out
}
