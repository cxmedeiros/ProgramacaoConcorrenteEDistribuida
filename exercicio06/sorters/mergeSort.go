package sorters

const mergeSortThreshold = 100 // Ajuste este valor conforme necessário

func MergeSort(arr []int64) []int64 {
	if len(arr) <= 1 {
		return arr
	}

	mid := len(arr) / 2
	left := MergeSort(arr[:mid])
	right := MergeSort(arr[mid:])

	return merge(left, right)
}

func MergeSortAsync(arr []int64) []int64 {
	done := make(chan []int64)

	go concurrentMergeSort(arr, done, 0)

	return <-done
}

func concurrentMergeSort(arr []int64, done chan []int64, depth int) {
	if len(arr) <= 1 {
		done <- arr
		return
	}

	mid := len(arr) / 2

	if depth < mergeSortThreshold {
		leftDone := make(chan []int64)
		rightDone := make(chan []int64)

		go concurrentMergeSort(arr[:mid], leftDone, depth+1)
		go concurrentMergeSort(arr[mid:], rightDone, depth+1)

		left := <-leftDone
		right := <-rightDone

		done <- merge(left, right)
	} else {
		left := MergeSort(arr[:mid])
		right := MergeSort(arr[mid:])

		done <- merge(left, right)
	}
}

func merge(left, right []int64) []int64 {
	result := make([]int64, 0, len(left)+len(right))
	i, j := 0, 0

	for i < len(left) && j < len(right) {
		if left[i] < right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}

	result = append(result, left[i:]...)
	result = append(result, right[j:]...)

	return result
}
