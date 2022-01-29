package array

func StringFilter(arr []string, filter func(string) bool) []string {
	var result []string
	for _, v := range arr {
		if filter(v) {
			result = append(result, v)
		}
	}
	return result
}
