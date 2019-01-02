package main

import (
	"fmt"
	"sort"
)

type IntSlice []int

func (s IntSlice) Len() int           { return len(s) }
func (s IntSlice) Less(i, j int) bool { return s[i] < s[j] }
func (s IntSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

var list1, list2, list3 IntSlice

func IsPalindrome(s sort.Interface) bool {
	for beg, end := 0, s.Len()-1; beg < end; beg, end = beg+1, end-1 {
		if !s.Less(beg, end) && !s.Less(end, beg) {
			continue
		} else {
			return false
		}
	}
	return true
}

func main() {
	fmt.Printf("list1:\n%v\n", list1)
	fmt.Printf("Is palindrome ? %v\n", IsPalindrome(list1))

	fmt.Printf("list2:\n%v\n", list2)
	fmt.Printf("Is palindrome ? %v\n", IsPalindrome(list2))

	fmt.Printf("list3:\n%v\n", list3)
	fmt.Printf("Is palindrome ? %v\n", IsPalindrome(list3))
}

func init() {
	list1 = IntSlice([]int{1, 2, 3, 3, 2, 1})
	list2 = IntSlice([]int{1, 2, 3, 4, 5})
	list3 = IntSlice([]int{1})
}
