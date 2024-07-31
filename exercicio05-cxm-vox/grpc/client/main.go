package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	pb "Vituriano/sort_grpc/sort_grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func measureTime(
	grpcMethod func(ctx context.Context, req *pb.SortRequest, opts ...grpc.CallOption) (*pb.SortResponse, error),
	req *pb.SortRequest,
) (*pb.SortResponse, time.Duration, error) {
	startTime := time.Now()
	res, err := grpcMethod(context.Background(), req)
	duration := time.Since(startTime)
	return res, duration, err
}

func sort(times int, mergeSortList []string, quickSortList []string) {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewSortServiceClient(conn)
	req := &pb.SortRequest{Array: []int32{38, 27, 43, 3, 9, 82, 10}}

	for i := 0; i < times; i++ {
		resp, duration, err := measureTime(client.MergeSort, req)
		if err != nil {
			log.Fatalf("could not call MergeSort: %v", err)
		}
		fmt.Println("Merge Sorted array:", resp.SortedArray)
		fmt.Println("MergeSort took", duration)
		mergeSortList[i] = duration.String()
	}

	for i := 0; i < times; i++ {
		resp, duration, err := measureTime(client.QuickSort, req)
		if err != nil {
			log.Fatalf("could not call QuickSort: %v", err)
		}
		fmt.Println("Quick Sorted array:", resp.SortedArray)
		fmt.Println("QuickSort took", duration)
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
