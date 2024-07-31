package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
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
	f, _ := os.Create("./canal.txt")
	defer f.Close()
	w := bufio.NewWriter(f)
	req := Request{
		Array:    []int{38, 27, 43, 3, 9, 82, 10},
		SortType: "merge",
	}

	fmt.Println("mergeSort")
	_, _ = fmt.Fprintf(w, "MergeSort\n")
	for i := 1; i < 1000; i++ {
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
		//fmt.Println("Sorted array:", resp.SortedArray)
		_, _ = fmt.Fprintf(w, "%v\n", duration)

		w.Flush()
		fmt.Println(duration)
	}

	fmt.Println("quickSort:")
	_, _ = fmt.Fprintf(w, "QuickSort\n")
	for i := 0; i < 1000; i++ {
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
		//fmt.Println("Sorted array:", resp.SortedArray)
		_, _ = fmt.Fprintf(w, "%v\n", duration)

		w.Flush()
		fmt.Println(duration)
	}
}
