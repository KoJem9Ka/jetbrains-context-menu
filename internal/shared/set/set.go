package set

// Set is a custom implementation of a set using a map.
type Set[T comparable] struct {
	items map[T]struct{}
}

// NewSet creates a new Set.
func NewSet[T comparable]() Set[T] {
	return Set[T]{
		items: make(map[T]struct{}),
	}
}

// Add adds an item to the set.
func (s *Set[T]) Add(item T) {
	s.items[item] = struct{}{}
}

// Remove removes an item from the set.
func (s *Set[T]) Remove(item T) {
	delete(s.items, item)
}

// Has checks if an item is in the set.
func (s *Set[T]) Has(item T) bool {
	_, exists := s.items[item]
	return exists
}

func (s *Set[T]) Clear() {
	s.items = make(map[T]struct{})
}

// Len returns the number of items in the set.
func (s *Set[T]) Len() int {
	return len(s.items)
}
