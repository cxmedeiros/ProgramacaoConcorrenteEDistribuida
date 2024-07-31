package main

import (
	"encoding/json"
	"fmt"
	"net"
)

type Request struct {
	Array    []int
	SortType string
}

type Response struct {
	SortedArray []int
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	var req Request
	err := decoder.Decode(&req)
	if err != nil {
		fmt.Println("Error decoding request:", err)
		return
	}

	var sortedArray []int
	if req.SortType == "merge" {
		sortedArray = MergeSort(req.Array)
	} else if req.SortType == "quick" {
		arrayCopy := make([]int, len(req.Array))
		copy(arrayCopy, req.Array)
		QuickSort(arrayCopy, 0, len(arrayCopy)-1)
		sortedArray = arrayCopy
	} else {
		fmt.Println("Invalid sort type")
		return
	}

	resp := Response{SortedArray: sortedArray}
	encoder := json.NewEncoder(conn)
	err = encoder.Encode(&resp)
	if err != nil {
		fmt.Println("Error encoding response:", err)
		return
	}
}

func main() {
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println("Error starting the server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server listening on port 1234...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}
