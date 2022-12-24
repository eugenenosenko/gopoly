package xmaps

func Merge[M map[K]V, K comparable, V any](maps ...M) M {
	res := make(M, 0)
	for _, m := range maps {
		for k, v := range m {
			res[k] = v
		}
	}
	return res
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
