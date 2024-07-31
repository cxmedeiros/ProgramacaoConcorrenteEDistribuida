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

func (s *server) MergeSortAsync(ctx context.Context, req *pb.SortRequest) (*pb.SortResponse, error) {
	sortedArray := sorters.MergeSortAsync(req.Array)
	return &pb.SortResponse{SortedArray: sortedArray}, nil
}

func (s *server) QuickSortAsync(ctx context.Context, req *pb.SortRequest) (*pb.SortResponse, error) {
	sorters.QuickSortAsync(req.Array, 0, int64(len(req.Array)-1))

	return &pb.SortResponse{SortedArray: req.Array}, nil
}

func (s *server) MergeSort(ctx context.Context, req *pb.SortRequest) (*pb.SortResponse, error) {
	sortedArray := sorters.MergeSort(req.Array)
	return &pb.SortResponse{SortedArray: sortedArray}, nil
}

func (s *server) QuickSort(ctx context.Context, req *pb.SortRequest) (*pb.SortResponse, error) {
	intArray := req.Array
	sorters.QuickSort(intArray, 0, int64(len(intArray)-1))
	sortedArray := intArray
	return &pb.SortResponse{SortedArray: sortedArray}, nil
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
