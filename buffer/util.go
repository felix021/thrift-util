package buffer

// IsPowerOf2 Check if the number is greater than 0 and has only one set bit (isolate the rightmost 1-bit)
func IsPowerOf2(n int) bool {
	return n > 0 && (n&(n-1)) == 0
}

func GetExponent(n int) (i int) {
	if !IsPowerOf2(n) {
		return -1
	}
	for ; n > 1; i++ {
		n /= 2
	}
	return
}
