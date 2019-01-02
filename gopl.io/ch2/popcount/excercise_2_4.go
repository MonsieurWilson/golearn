package popcount

// PopCountByShifting returns the population count (number of set bits) of x.
func PopCountByShifting(x uint64) int {
	var count int
	for x != 0 {
		count += int(x & 1)
		x >>= 1
	}
	return count
}
