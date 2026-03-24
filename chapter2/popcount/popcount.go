package popcount

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
