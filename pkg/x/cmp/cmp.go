package cmp

// Or returns the first of its arguments that is not equal to the zero value.
// If no argument is non-zero, it returns the zero value.
func Or[T comparable](vals ...T) T {
	var zero T
	for _, val := range vals {
		if val != zero {
			return val
		}
	}
	return zero
}

func IfElse[T any](conditon bool, v1, v2 T) T {
	if conditon {
		return v1
	}
	return v2
}
