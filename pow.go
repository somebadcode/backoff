package backoff

type integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

func pow[T integer](a, b T) T {
	var result T = 1

	for b != 0 {
		if b&1 != 0 {
			result *= a
		}

		b >>= 1
		a *= a
	}

	return result
}
