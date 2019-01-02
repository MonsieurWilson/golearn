package popcount

// PopCountByClearing returns the population count (number of set bits) of x.
func PopCountByClearing(x uint64) int {
	var count int
	for x != 0 {
		count += 1
		x = x & (x - 1)
	}
	return count
}
