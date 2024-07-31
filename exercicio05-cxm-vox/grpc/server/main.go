package main

import (
	"context"
	"log"
	"net"

	pb "Vituriano/sort_grpc/sort_grpc"
	"Vituriano/sorters"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedSortServiceServer
}

func (s *server) MergeSort(ctx context.Context, req *pb.SortRequest) (*pb.SortResponse, error) {
	sortedArray := mergeSort(req.Array)
	return &pb.SortResponse{SortedArray: sortedArray}, nil
}

func (s *server) QuickSort(ctx context.Context, req *pb.SortRequest) (*pb.SortResponse, error) {
	sortedArray := quickSort(req.Array)
	return &pb.SortResponse{SortedArray: sortedArray}, nil
}

func mergeSort(arr []int32) []int32 {
	// Converter []int32 para []int
	intArray := convert(arr, func(v int32) int {
		return int(v)
	})

	intArray = sorters.MergeSort(intArray)

	// Converter []int de volta para []int32
	int32Array := convert(intArray, func(v int) int32 {
		return int32(v)
	})

	return int32Array
}

func quickSort(arr []int32) []int32 {
	// Converter []int32 para []int
	intArray := convert(arr, func(v int32) int {
		return int(v)
	})

	sorters.QuickSort(intArray, 0, len(arr)-1)

	// Converter []int de volta para []int32
	int32Array := convert(intArray, func(v int) int32 {
		return int32(v)
	})

	return int32Array
}

func convert[T, U any](arr []T, convertFunc func(T) U) []U {
	result := make([]U, len(arr))
	for i, v := range arr {
		result[i] = convertFunc(v)
	}
	return result
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterSortServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
