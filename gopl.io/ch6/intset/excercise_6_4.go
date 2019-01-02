// Package intset provides a set of integers based on a bit vector.
package intset

func (s *IntSet) Elems() []uint64 {
	var elements []uint64
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < 64; j++ {
			if word&(1<<uint64(j)) != 0 {
				elements = append(elements, uint64(64*i+j))
			}
		}
	}
	return elements
}
