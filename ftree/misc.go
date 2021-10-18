package ftree

// numDigits returns the number of digits in the integer if it was
// formatted as a string.
func numDigits(i uint64) (n int) {
	for i > 9 {
		i /= 10
		n++
	}

	return n + 1
}
