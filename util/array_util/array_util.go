package array_util

func IntMinMax(array []int) (int, int) {
	var max int = array[0]
	var min int = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}
func IntMin(array []int) int {
	min := array[0]
	for _, value := range array {
		if min > value {
			min = value
		}
	}
	return min
}
