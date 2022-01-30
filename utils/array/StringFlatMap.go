package array

func StringFlatMap(arr []string, f func(string) []string) []string {
	result := make([]string, 0)
	for _, v := range arr {
		result = append(result, f(v)...)
	}
	return result
}
