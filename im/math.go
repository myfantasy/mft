// Package im - int math simple operations lake max min etc for int
package im

// Max - maximum of 2 elements
func Max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

// MaxS - maximum of Slice elements; if slice is empty then 0
func MaxS(s ...int) int {
	if len(s) == 0 {
		return 0
	}

	m := s[0]
	for k := 1; k < len(s); k++ {
		if s[k] > m {
			m = s[k]
		}
	}
	return m
}

// Min - minimum of 2 elements
func Min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

// MinS - minimum of Slice elements; if slice is empty then 0
func MinS(s ...int) int {
	if len(s) == 0 {
		return 0
	}

	m := s[0]
	for k := 1; k < len(s); k++ {
		if s[k] < m {
			m = s[k]
		}
	}
	return m
}
