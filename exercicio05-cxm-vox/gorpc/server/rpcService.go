package main

import (
	"fmt"
	"net"
	"net/rpc"

	"Vituriano/sorters"
)

type SortService struct{}

type Args struct {
	Array []int
}

type SortResponse struct {
	SortedArray []int
}

func (s *SortService) MergeSortRemote(args *Args, reply *SortResponse) error {
	reply.SortedArray = sorters.MergeSort(args.Array)
	return nil
}

func (s *SortService) QuickSortRemote(args *Args, reply *SortResponse) error {
	array := make([]int, len(args.Array))
	copy(array, args.Array)
	sorters.QuickSort(array, 0, len(array)-1)
	reply.SortedArray = array
	return nil
}

func main() {
	sortService := new(SortService)
	rpc.Register(sortService)

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println("Error starting the server:", err)
		return
	}

	fmt.Println("Server listening on port 1234...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
