package service

import "unicode/utf8"

func StringLenMin(s string, min int) bool {
	if len(s) < min {
		return false
	}

	if utf8.RuneCountInString(s) < min {
		return false
	}

	return true
}

func StringLenMax(s string, max int) bool {
	if len(s)/utf8.MaxRune > max {
		return false
	}

	if utf8.RuneCountInString(s) > max {
		return false
	}

	return true
}

func StringLenBetween(s string, min, max int) int {
	if len(s) < min {
		return -1
	}

	if len(s)/utf8.MaxRune > max {
		return 1
	}

	count := utf8.RuneCountInString(s)

	if count < min {
		return -1
	}

	if count > max {
		return 1
	}

	return 0
}
