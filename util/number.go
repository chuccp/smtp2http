package util

func ContainsNumberAny(s uint, strs ...uint) bool {
	for _, str := range strs {
		if s == str {
			return true
		}
	}
	return false
}
