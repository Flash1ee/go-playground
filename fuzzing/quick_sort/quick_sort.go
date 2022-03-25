package quick_sort

func partition(array []int) int {
	pivot := array[len(array)/2]
	left := 0
	right := len(array) - 1
	for {

		for ; array[left] < pivot; left++ {
		}
		for ; array[right] > pivot; right-- {
		}

		if left >= right || array[left] == array[right] {
			return right
		}
		array[left], array[right] = array[right], array[left]
	}
}

func QuickSort(array []int) []int {
	if len(array) < 2 {
		return array
	}
	split := partition(array)
	QuickSort(array[:split])
	QuickSort(array[split:])

	return array
}
