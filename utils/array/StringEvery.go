package array

func StringEvery(array []string, predicate func(string) bool) bool {
	for _, v := range array {
		if !predicate(v) {
			return false
		}
	}
	return true
}
