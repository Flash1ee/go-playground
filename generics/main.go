package main

// написать функцию возвращающую слайс уникальных значений из входящего слайса произвольного типа содержащего повторы

func getUniqNumSlice(src []int) (dst []int) {
	dst = make([]int, 0)
	exists := make(map[int]struct{})

	for _, val := range src {
		if _, ok := exists[val]; !ok {
			dst = append(dst, val)
			exists[val] = struct{}{}
		}
	}
	return
}

func getTemplateUniq[T any](src []T) (dst []T) {
	dst = make([]T, 0)
	exists := make(map[any]struct{})

	for _, val := range src {
		if _, ok := exists[val]; !ok {
			dst = append(dst, val)
			exists[val] = struct{}{}
		}
	}
	return
}

func main() {
}
