package main

type number interface {
	float64 | float32 | int64 | int32 | int16 | int8 | int | uint64 | uint32 | uint16 | uint8 | uint
}

func min[T number](a T, bc ...T) T {
	m := a
	for _, n := range bc {
		if n < m {
			m = n
		}
	}
	return m
}

func max[T number](a T, bc ...T) T {
	m := a
	for _, n := range bc {
		if n > m {
			m = n
		}
	}
	return m
}
