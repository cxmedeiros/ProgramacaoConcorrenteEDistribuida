package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

func multiplyMatrices(A, B [][]int) [][]int {
	rowsA := len(A)
	colsA := len(A[0])
	rowsB := len(B)
	colsB := len(B[0])

	if colsA != rowsB {
		panic("Número de colunas de A não é igual ao número de linhas de B")
	}

	C := make([][]int, rowsA)
	for i := range C {
		C[i] = make([]int, colsB)
	}

	for i := 0; i < rowsA; i++ {
		for j := 0; j < colsB; j++ {
			for k := 0; k < colsA; k++ {
				C[i][j] += A[i][k] * B[k][j]
			}
		}
	}
	return C
}

func multiplyRowColumn(A, B, C [][]int, row, col int, wg *sync.WaitGroup, mtx *sync.Mutex) {
	defer wg.Done()
	sum := 0
	for i := range A[row] {
		sum += A[row][i] * B[i][col]
	}
	mtx.Lock()
	C[row][col] = sum
	mtx.Unlock()
}

func multiplyMatricesConcurrent(A, B [][]int) [][]int {
	rowsA := len(A)
	colsA := len(A[0])
	rowsB := len(B)
	colsB := len(B[0])

	if colsA != rowsB {
		panic("Número de colunas de A não é igual ao número de linhas de B")
	}

	C := make([][]int, rowsA)
	for i := range C {
		C[i] = make([]int, colsB)
	}

	var wg sync.WaitGroup
	var mtx sync.Mutex

	for i := 0; i < rowsA; i++ {
		for j := 0; j < colsB; j++ {
			wg.Add(1)
			go multiplyRowColumn(A, B, C, i, j, &wg, &mtx)
		}
	}

	wg.Wait()
	return C
}

func main() {
	f, _ := os.Create("./results.txt")
	defer f.Close()
	w := bufio.NewWriter(f)

	sizes := []int{10, 100, 1000}

	for _, size := range sizes {
		A := generateRandomMatrix(size, size)
		B := generateRandomMatrix(size, size)

		fmt.Printf("Testando com matrizes de tamanho %dx%d\n", size, size)
		_, _ = fmt.Fprintf(w, "Testando com matrizes de tamanho %dx%d\n", size, size)
		fmt.Printf("Sem concorrência:\n")
		_, _ = fmt.Fprintf(w, "Sem concorrência:\n")
		w.Flush()

		for i := 0; i < 30; i++ {
			start := time.Now()
			multiplyMatrices(A, B)
			elapsed := time.Since(start)
			fmt.Printf(" %v\n", elapsed)
			_, _ = fmt.Fprintf(w, "%v\n", elapsed)
			w.Flush()
		}

		fmt.Printf("Com concorrência:\n")
		_, _ = fmt.Fprintf(w, "Com concorrência:\n")
		for i := 0; i < 30; i++ {
			start := time.Now()
			multiplyMatricesConcurrent(A, B)
			elapsed := time.Since(start)
			fmt.Printf("%v\n", elapsed)
			_, _ = fmt.Fprintf(w, "%v\n", elapsed)
			w.Flush()
		}
	}
	w.Flush()
}

func generateRandomMatrix(rows, cols int) [][]int {
	matrix := make([][]int, rows)
	for i := range matrix {
		matrix[i] = make([]int, cols)
		for j := range matrix[i] {
			matrix[i][j] = rand.Intn(100)
		}
	}
	return matrix
}
