package utils

func LastIndex[S ~[]E, E comparable](s S, v E) int {
	for i := len(s) - 1; i >= 0; i-- {
		if v == s[i] {
			return i
		}
	}
	return -1
}

func LastIndexFunc[S ~[]E, E any](s S, f func(E) bool) int {
	for i := len(s) - 1; i >= 0; i-- {
		if f(s[i]) {
			return i
		}
	}
	return -1
}

func FilterFunc[S ~[]E, E any](s S, f func(E) bool) S {
	r := make(S, 0, len(s))

	for _, v := range s {
		if f(v) {
			r = append(r, v)
		}
	}

	return r
}

func MapFunc[S ~[]E, E any, R any](s S, f func(E) R) []R {
	r := make([]R, len(s))

	for i, v := range s {
		r[i] = f(v)
	}

	return r
}
