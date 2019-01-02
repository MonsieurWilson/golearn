// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 165.

// Package intset provides a set of integers based on a bit vector.
package intset2

import (
	"bytes"
	"fmt"
)

//!+intset

// An IntSet is a set of small non-negative integers.
// Its zero value represents the empty set.
type IntSet struct {
	words []uint
}

const OpMode = 32 << (^uint(0) >> 63)

// Has reports whether the set contains the non-negative value x.
func (s *IntSet) Has(x int) bool {
	word, bit := x/OpMode, uint(x%OpMode)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

// Add adds the non-negative value x to the set.
func (s *IntSet) Add(x int) {
	word, bit := x/OpMode, uint(x%OpMode)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

// UnionWith sets s to the union of s and t.
func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

//!-intset

//!+string

// String returns the set as a string of the form "{1 2 3}".
func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < OpMode; j++ {
			if word&(1<<uint(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", OpMode*i+j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

//!-string

// Calculate the number of elements.
func (s *IntSet) Len() (num int) {
	for _, word := range s.words {
		num += bitcount(word)
	}
	return
}

func bitcount(bit uint) int {
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
		word, bit := x/OpMode, uint(x%OpMode)
		s.words[word] &^= 1 << bit
	}
}

// Clear set.
func (s *IntSet) Clear() {
	s.words = []uint{}
}

// Copy set.
func (s *IntSet) Copy() *IntSet {
	copyset := IntSet{make([]uint, s.Len())}
	copy(copyset.words, s.words)
	return &copyset
}

// Add adds all the given non-negative values to the set.
func (s *IntSet) AddAll(vals ...int) {
	for _, val := range vals {
		s.Add(val)
	}
}

func (s *IntSet) Elems() []int {
	var elements []int
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < OpMode; j++ {
			if word&(1<<uint(j)) != 0 {
				elements = append(elements, OpMode*i+j)
			}
		}
	}
	return elements
}
