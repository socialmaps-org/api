package j

func Get[T any](v any, path ...any) (r T) {
	cur := v
	for _, p := range path {
		switch key := p.(type) {
		case string:
			m, ok := cur.(map[string]any)
			if !ok {
				panic("not map")
			}
			cur, ok = m[key]
			if !ok {
				panic("key not found")
			}
		case int:
			a, ok := cur.([]any)
			if !ok {
				panic("not array")
			}
			if key < 0 || key >= len(a) {
				panic("invalid index")
			}
			cur = a[key]
		default:
			panic("invalid path")
		}
	}

	curt, ok := cur.(T)
	if ok {
		return curt
	}

	f64, ok := cur.(float64)
	if !ok {
		panic("invalid type")
	}

	switch any(r).(type) {
	case int64:
		return any(int64(f64)).(T)
	case int32:
		return any(int32(f64)).(T)
	case int:
		return any(int(f64)).(T)
	default:
		panic("unsupported type")
	}
}
