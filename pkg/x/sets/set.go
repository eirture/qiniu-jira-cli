package sets

type Set[T comparable] map[T]struct{}

func NewSet[T comparable](items ...T) Set[T] {
	return make(Set[T]).Add(items...)
}

func (s Set[T]) Add(items ...T) Set[T] {
	for _, it := range items {
		s[it] = struct{}{}
	}
	return s
}

func (s Set[T]) Del(items ...T) Set[T] {
	for _, it := range items {
		delete(s, it)
	}
	return s
}

func (s Set[T]) Has(item T) bool {
	_, ok := s[item]
	return ok
}
