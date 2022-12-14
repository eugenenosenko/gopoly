package xslices

import (
	"fmt"
)

func Map[S ~[]T, O ~[]R, T, R any](s S, mapper func(T) R) O {
	out := make([]R, 0, len(s))
	for _, t := range s {
		out = append(out, mapper(t))
	}
	return out
}

func Flatten[S ~[][]T, T any](s S) []T {
	target := make([]T, 0)
	for _, ts := range s {
		target = append(target, ts...)
	}
	return target
}

func ToSet[S ~[]T, T comparable](s S) map[T]struct{} {
	return ToSetFunc[S](s, func(t T) T {
		return t
	})
}

func ToSetFunc[S ~[]T, T, V comparable](s S, kMapper func(T) V) map[V]struct{} {
	return ToMap[S, map[V]struct{}](s, kMapper, func(_ T) struct{} { return struct{}{} })
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
			r, ok := any(t).(R)
			if !ok {
				panic(fmt.Sprintf("cant cast %T", t))
			}
			out[k] = r
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
