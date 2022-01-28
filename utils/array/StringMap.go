package array

func StringMap(arr []string, fn func(str string) string) []string {
	for i, s := range arr {
		arr[i] = fn(s)
	}
	return arr
}
