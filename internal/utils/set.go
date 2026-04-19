package utils

type Set[T comparable] struct {
	kv map[T]struct{}
}

func NewSet[T comparable](values ...T) *Set[T] {
	kv := make(map[T]struct{}, len(values))
	for _, value := range values {
		kv[value] = struct{}{}
	}
	return &Set[T]{kv}
}

func (s *Set[T]) Has(value T) bool {
	_, ok := s.kv[value]
	return ok
}
