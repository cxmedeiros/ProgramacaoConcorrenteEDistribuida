package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type Request struct {
	Array    []int
	SortType string
}

type Response struct {
	SortedArray []int
}

func main() {
	req := Request{
		Array:    []int{38, 27, 43, 3, 9, 82, 10},
		SortType: "merge",
	}

	for i := 1; i < 10; i++ {
		conn, err := net.Dial("tcp", "localhost:1234")
		if err != nil {
			fmt.Println("Error dialing the server:", err)
			return
		}
		defer conn.Close()

		start := time.Now()
		encoder := json.NewEncoder(conn)
		err = encoder.Encode(&req)
		if err != nil {
			fmt.Println("Error encoding request:", err)
			return
		}

		decoder := json.NewDecoder(conn)
		var resp Response
		err = decoder.Decode(&resp)
		if err != nil {
			fmt.Println("Error decoding response:", err)
			return
		}
		duration := time.Since(start)
		fmt.Println("Sorted array:", resp.SortedArray)
		fmt.Println("MergeSort took", duration)
	}

	for i := 0; i < 10; i++ {
		conn, err := net.Dial("tcp", "localhost:1234")
		if err != nil {
			fmt.Println("Error dialing the server:", err)
			return
		}
		defer conn.Close()

		req.SortType = "quick"
		start := time.Now()
		encoder := json.NewEncoder(conn)
		err = encoder.Encode(&req)
		if err != nil {
			fmt.Println("Error encoding request:", err)
			return
		}

		decoder := json.NewDecoder(conn)
		var resp Response
		err = decoder.Decode(&resp)
		if err != nil {
			fmt.Println("Error decoding response:", err)
			return
		}
		duration := time.Since(start)
		fmt.Println("Sorted array:", resp.SortedArray)
		fmt.Println("QuickSort took", duration)
	}
}
