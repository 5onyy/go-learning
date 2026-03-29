package popcount

/* IMPPORTANT NOTES: visibility of the function to be use outside package is by the Capital of the first letter*/

func PopCount(x uint64) int {
	var ans int = 0
	for ; x > 0; x /= 2 {
		ans += int(x % 2)
	}
	return ans
}

func PopCount2(x uint64) int {
	var ans int = 0
	for ; x > 0; x = x & (x - 1) {
		ans++
	}
	return ans
}
