// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 101.

// Package tree provides insertion sort using an unbalanced binary tree.
package tree

//!+
type BinaryTree struct {
	value       int
	left, right *BinaryTree
}

// to be continued
func (t *BinaryTree) String() string {

}

// Sort sorts values in place.
func Sort(values []int) {
	var root *BinaryTree
	for _, v := range values {
		root = add(root, v)
	}
	appendValues(values[:0], root)
}

// appendValues appends the elements of t to values in order
// and returns the resulting slice.
func appendValues(values []int, t *BinaryTree) []int {
	if t != nil {
		values = appendValues(values, t.left)
		values = append(values, t.value)
		values = appendValues(values, t.right)
	}
	return values
}

func add(t *BinaryTree, value int) *BinaryTree {
	if t == nil {
		// Equivalent to return &BinaryTree{value: value}.
		t = new(BinaryTree)
		t.value = value
		return t
	}
	if value < t.value {
		t.left = add(t.left, value)
	} else {
		t.right = add(t.right, value)
	}
	return t
}

//!-
