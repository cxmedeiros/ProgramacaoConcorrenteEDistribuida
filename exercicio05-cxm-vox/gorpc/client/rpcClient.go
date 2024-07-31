package main

import (
	"encoding/csv"
	"fmt"
	"net/rpc"
	"os"
	"time"
)

type Args struct {
	Array []int
}

type SortResponse struct {
	SortedArray []int
}

func measureTime(client *rpc.Client, method string, args *Args) (SortResponse, time.Duration, error) {
	var reply SortResponse
	start := time.Now()
	err := client.Call(method, args, &reply)
	duration := time.Since(start)
	return reply, duration, err
}

func sort(times int, mergeSortList []string, quickSortList []string) {
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		fmt.Println("Error dialing the server:", err)
		return
	}
	defer client.Close()

	args := &Args{Array: []int{38, 27, 43, 3, 9, 82, 10}}
	for i := 0; i < times; i++ {
		reply, duration, err := measureTime(client, "SortService.MergeSortRemote", args)
		if err != nil {
			fmt.Println("Error calling MergeSortRemote:", err)
			return
		}
		fmt.Println("Merge Sorted array:", reply.SortedArray)
		fmt.Println("MergeSortRemote took", duration)
		mergeSortList[i] = duration.String()

	}

	for i := 0; i < times; i++ {
		reply, duration, err := measureTime(client, "SortService.QuickSortRemote", args)
		if err != nil {
			fmt.Println("Error calling QuickSortRemote:", err)
			return
		}
		fmt.Println("Quick Sorted array:", reply.SortedArray)
		fmt.Println("QuickSortRemote took", duration)
		quickSortList[i] = duration.String()
	}
}

func main() {
	timestamp := time.Now().Format("2006-01-02_15h-04m-05s")
	filePath := "../temp/" + timestamp + ".csv"

	f, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer f.Close()
	w := csv.NewWriter(f)
	header := []string{
		"Merge Sort - 100", "Quick Sort - 100",
		"Merge Sort - 1000", "Quick Sort - 1000",
		"Merge Sort - 10000", "Quick Sort - 10000",
	}

	mergeSortList100 := make([]string, 100)
	quickSortList100 := make([]string, 100)
	mergeSortList1000 := make([]string, 1000)
	quickSortList1000 := make([]string, 1000)
	mergeSortList10000 := make([]string, 10000)
	quickSortList10000 := make([]string, 10000)

	sort(100, mergeSortList100, quickSortList100)
	sort(1000, mergeSortList1000, quickSortList1000)
	sort(10000, mergeSortList10000, quickSortList10000)

	getValue := func(list []string, i int) string {
		if i < len(list) {
			return list[i]
		}
		return ""
	}

	data := [][]string{}
	maxLen := 10000
	for i := 0; i < maxLen; i++ {
		row := []string{
			getValue(mergeSortList100, i),
			getValue(quickSortList100, i),
			getValue(mergeSortList1000, i),
			getValue(quickSortList1000, i),
			getValue(mergeSortList10000, i),
			getValue(quickSortList10000, i),
		}
		data = append(data, row)
	}

	w.Write(header)
	for _, row := range data {
		w.Write(row)
	}

	w.Flush()
}
