package object

func Assign(data map[string]string, others ...map[string]string) map[string]string {
	for _, other := range others {
		for k, v := range other {
			data[k] = v
		}
	}
	return data
}
