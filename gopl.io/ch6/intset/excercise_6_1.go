// Package intset provides a set of integers based on a bit vector.
package intset

// Calculate the number of elements.
func (s *IntSet) Len() (num int) {
	for _, word := range s.words {
		num += bitcount(word)
	}
	return
}

func bitcount(bit uint64) int {
	var nbits int
	for bit != 0 {
		bit = bit & (bit - 1) // clear rightmost non-zero bit
		nbits++
	}
	return nbits
}

// Remove element from set.
func (s *IntSet) Remove(x int) {
	if s.Has(x) {
		word, bit := x/64, uint(x%64)
		s.words[word] &^= 1 << bit
	}
}

// Clear set.
func (s *IntSet) Clear() {
	s.words = []uint64{}
}

// Copy set.
func (s *IntSet) Copy() *IntSet {
	copyset := IntSet{make([]uint64, s.Len())}
	copy(copyset.words, s.words)
	return &copyset
}
