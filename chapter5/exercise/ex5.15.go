package main

func find_min(required int, optional ...int) int {
	mn := required
	for _, val := range optional {
		mn = min(mn, val)
	}
	return mn
}

func find_max(required int, optional ...int) int {
	mx := required
	for _, val := range optional {
		mx = max(mx, val)
	}
	return mx
}

func main() {

}
