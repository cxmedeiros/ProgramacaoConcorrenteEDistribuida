package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type CellResult struct {
	row int
	col int
	val int
}

func multiplyRowColumn(A, B [][]int, row, col int, resultChan chan<- CellResult) {
	sum := 0
	for i := range A[row] {
		sum += A[row][i] * B[i][col]
	}
	resultChan <- CellResult{row: row, col: col, val: sum}
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

	resultChan := make(chan CellResult, rowsA*colsB)

	for i := 0; i < rowsA; i++ {
		for j := 0; j < colsB; j++ {
			go multiplyRowColumn(A, B, i, j, resultChan)
		}
	}

	for i := 0; i < rowsA*colsB; i++ {
		result := <-resultChan
		C[result.row][result.col] = result.val
	}

	close(resultChan)

	return C
}

func main() {
	f, _ := os.Create("./canal.txt")
	defer f.Close()
	w := bufio.NewWriter(f)

	sizes := []int{10, 100, 1000}

	for _, size := range sizes {
		A := generateRandomMatrix(size, size)
		B := generateRandomMatrix(size, size)

		fmt.Printf("Testando com matrizes de tamanho %dx%d\n", size, size)
		_, _ = fmt.Fprintf(w, "Testando com matrizes de tamanho %dx%d\n", size, size)
		w.Flush()

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
