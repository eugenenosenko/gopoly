package xmap

func KeysFunc[M ~map[K]V, O ~[]T, K, V comparable, T any](m M, f func(K) T) O {
	out := make(O, 0, len(m))
	for k, _ := range m {
		out = append(out, f(k))
	}
	return out
}

func Difference[M map[K]V, K, V comparable](m1, m2 M) M {
	diff := make(map[K]V)
	for k, v := range m1 {
		if e, ok := m2[k]; ok {
			if e != v {
				diff[k] = e
			}
		} else {
			diff[k] = v
		}
	}

	for k, v := range m2 {
		if _, ok := m1[k]; !ok {
			diff[k] = v
		}
	}
	return diff
}
