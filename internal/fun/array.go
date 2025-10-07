package fun

func Filter[T any](xs []T, cond func(x T) bool) []T {
	res := []T{}
	for _, x := range xs {
		if cond(x) {
			res = append(res, x)
		}
	}
	return res
}
