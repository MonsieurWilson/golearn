// Package intset provides a set of integers based on a bit vector.
package intset

// Add adds all the given non-negative values to the set.
func (s *IntSet) AddAll(vals ...int) {
	for _, val := range vals {
		s.Add(val)
	}
}
