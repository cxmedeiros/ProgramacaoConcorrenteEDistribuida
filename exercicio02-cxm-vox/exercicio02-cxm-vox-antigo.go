package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

func MultiplicaMatrizesAsync(A, B [][]int) ([][]int, error) {
	numLinhasA := len(A)
	numColunasA := len(A[0])
	numLinhasB := len(B)
	numColunasB := len(B[0])

	if numColunasA != numLinhasB {
		return nil, fmt.Errorf("as matrizes não podem ser multiplicadas")
	}

	C := make([][]int, numLinhasA)
	for i := range C {
		C[i] = make([]int, numColunasB)
	}

	for i := 0; i < numLinhasA; i++ {
		for j := 0; j < numColunasB; j++ {
			for k := 0; k < numColunasA; k++ {
				C[i][j] += A[i][k] * B[k][j]
			}
		}
	}

	return C, nil
}

func MultiplicaMatrizesSync(A, B [][]int) ([][]int, error) {
	numLinhasA := len(A)
	numColunasA := len(A[0])
	numLinhasB := len(B)
	numColunasB := len(B[0])

	if numColunasA != numLinhasB {
		return nil, fmt.Errorf("as matrizes não podem ser multiplicadas")
	}

	C := make([][]int, numLinhasA)
	for i := range C {
		C[i] = make([]int, numColunasB)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 0; i < numLinhasA; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < numColunasB; j++ {
				for k := 0; k < numColunasA; k++ {
					mu.Lock()
					C[i][j] += A[i][k] * B[k][j]
					mu.Unlock()
				}
			}
		}(i)
	}

	wg.Wait()
	return C, nil
}

func GenerateMatrix(size int) [][]int {
	matrix := make([][]int, size)
	for i := range matrix {
		matrix[i] = make([]int, size)
		for j := range matrix[i] {
			matrix[i][j] = rand.Intn(100)
		}
	}
	return matrix
}

func main() {

	f, _ := os.Create("./results2.txt")
	defer f.Close()

	w := bufio.NewWriter(f)

	sizes := []int{10, 100, 1000}

	for _, size := range sizes {
		A := GenerateMatrix(size)
		B := GenerateMatrix(size)

		// fmt.Printf("Testando com matrizes de tamanho %dx%d\n", size, size)
		_, _ = fmt.Fprintf(w, "Testando com matrizes de tamanho %dx%d\n", size, size)
		w.Flush()

		for range 30 {
			start := time.Now()
			MultiplicaMatrizesAsync(A, B)
			// fmt.Printf("Async: %v\n", time.Since(start))
			finalTime := time.Since(start)
			_, _ = fmt.Fprintf(w, "Async, %v\n", finalTime)
			w.Flush()
		}

		for range 30 {
			start := time.Now()
			MultiplicaMatrizesSync(A, B)
			// fmt.Printf("Sync: %v\n", time.Since(start))
			finalTime := time.Since(start)
			_, _ = fmt.Fprintf(w, "Sync, %v\n", finalTime)
			w.Flush()
			// fmt.Printf("--------\n")
		}
	}
	w.Flush()
}
