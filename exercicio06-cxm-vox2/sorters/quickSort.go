package sorters

const quickSortThreshold = 100 // Ajuste este valor conforme necessário

func QuickSort(arr []int64, low, high int64) {
	if low < high {
		pi := partition(arr, low, high)
		QuickSort(arr, low, pi-1)
		QuickSort(arr, pi+1, high)
	}
}

func QuickSortAsync(arr []int64, low, high int64) []int64 {
	done := make(chan []int64)

	go concurrentQuickSort(arr, 0, int64(len(arr)-1), done, 0)

	return <-done
}

func concurrentQuickSort(arr []int64, low, high int64, done chan []int64, depth int) {
	if low < high {
		p := partition(arr, low, high)

		if depth < quickSortThreshold {
			leftDone := make(chan []int64)
			rightDone := make(chan []int64)

			go concurrentQuickSort(arr, low, p-1, leftDone, depth+1)
			go concurrentQuickSort(arr, p+1, high, rightDone, depth+1)

			<-leftDone
			<-rightDone
		} else {
			QuickSort(arr, low, p-1)
			QuickSort(arr, p+1, high)
		}
	}
	done <- arr
}

func partition(arr []int64, low, high int64) int64 {
	pivot := arr[high]
	i := low - 1

	for j := low; j < high; j++ {
		if arr[j] < pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}

	arr[i+1], arr[high] = arr[high], arr[i+1]
	return i + 1
}
