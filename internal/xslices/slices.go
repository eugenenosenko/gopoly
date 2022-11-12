package xslices

func Map[S ~[]T, O ~[]R, T, R any](s S, mapper func(T) R) O {
	out := make([]R, 0, len(s))
	for _, t := range s {
		out = append(out, mapper(t))
	}
	return out
}

func Filter[S ~[]T, T any](s S, predicate func(T) bool) S {
	out := make([]T, 0, len(s))
	for _, t := range s {
		if predicate(t) {
			out = append(out, t)
		}
	}
	return out
}

func Distinct[S ~[]T, T comparable](s S) S {
	m := make(map[T]struct{})
	res := make([]T, 0, len(s))
	for _, t := range s {
		if _, exist := m[t]; !exist {
			m[t] = struct{}{}
			res = append(res, t)
		}
	}
	return res
}

func ToMap[S ~[]T, O ~map[V]R, T, R, V comparable](s S, kMapper func(T) V, vMapper func(T) R) O {
	out := make(map[V]R, 0)
	for _, t := range s {
		if kMapper == nil {
			panic("key-mapper is nil")
		}
		k := kMapper(t)
		if vMapper != nil {
			out[k] = vMapper(t)
		} else {
			out[k] = any(t).(R)
		}
	}
	return out
}

func First[S ~[]E, E any](s S) (E, bool) {
	if len(s) == 0 {
		var zero E
		return zero, false
	}
	return s[0], true
}

func Difference[S ~[]T, T comparable](s1, s2 S) S {
	set := make(map[T]struct{}, len(s1))
	for _, t := range s1 {
		set[t] = struct{}{}
	}
	var diff []T
	for _, x := range s2 {
		if _, found := set[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}
