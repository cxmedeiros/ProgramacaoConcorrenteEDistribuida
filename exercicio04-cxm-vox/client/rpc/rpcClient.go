package main

import (
	"bufio"
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

func main() {
	f, _ := os.Create("./canal.txt")
	defer f.Close()
	w := bufio.NewWriter(f)
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		fmt.Println("Error dialing the server:", err)
		return
	}
	defer client.Close()

	args := &Args{Array: []int{38, 27, 43, 3, 9, 82, 10}}

	fmt.Println("MergeSort")
	_, _ = fmt.Fprintf(w, "MergeSort\n")

	for i := 0; i < 1000; i++ {
		_, duration, err := measureTime(client, "SortService.MergeSortRemote", args)
		if err != nil {
			fmt.Println("Error calling MergeSortRemote:", err)
			return
		}
		//fmt.Println("Merge Sorted array:", reply.SortedArray)
		_, _ = fmt.Fprintf(w, "%v\n", duration)
		w.Flush()
		fmt.Println(duration)
	}

	fmt.Println("QuickSort")
	_, _ = fmt.Fprintf(w, "QuickSort\n")
	for i := 0; i < 1000; i++ {
		_, duration, err := measureTime(client, "SortService.QuickSortRemote", args)
		if err != nil {
			fmt.Println("Error calling QuickSortRemote:", err)
			return
		}
		//fmt.Println("Quick Sorted array:", reply.SortedArray)
		_, _ = fmt.Fprintf(w, "%v\n", duration)
		w.Flush()
		fmt.Println(duration)
	}
}
