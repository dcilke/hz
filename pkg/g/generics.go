package g

func Includes[T comparable](slice []T, value T) bool {
	for _, e := range slice {
		if e == value {
			return true
		}
	}
	return false
}

func SameOrEmpty[T comparable](args ...T) (bool, T) {
	zero := *new(T)
	vals := Filter(args, func(i T) bool {
		return i != zero
	})
	if len(vals) == 0 {
		return true, zero
	}
	// if all args are all the same, return the value
	isSame := true
	for i := 1; i < len(vals); i++ {
		if vals[i-1] != vals[i] {
			isSame = false
			break
		}
	}
	if isSame {
		return true, vals[0]
	}
	// if there are multiple values, return the zero value
	return false, zero
}

func Filter[T comparable](slice []T, filter func(T) bool) []T {
	var result []T
	for _, e := range slice {
		if filter(e) {
			result = append(result, e)
		}
	}
	return result
}
