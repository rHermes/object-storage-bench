package ftree

// numDigits returns the number of digits in the integer if it was
// formatted as a string.
func numDigits(i uint64) int {
	n := 1
	for i > 9 {
		n++
		i /= 10
	}
	return n
}
