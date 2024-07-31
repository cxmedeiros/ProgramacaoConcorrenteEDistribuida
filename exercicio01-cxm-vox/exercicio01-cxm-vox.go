// Camila Xavier (CXM) & Vituriano Xisto (VOX)

package main

import (
	"fmt"
	"log"
	"sync"
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

func main() {
	A := [][]int{
		{3, 2},
		{5, -1},
	}

	B := [][]int{
		{6, 4, -2},
		{0, 7, 1},
	}

	fmt.Println("Async")
	C, err := MultiplicaMatrizesAsync(A, B)
	if err != nil {
		log.Fatal(err)
	}

	for _, linha := range C {
		fmt.Println(linha)
	}

	fmt.Println("\nSync")
	C, err = MultiplicaMatrizesSync(A, B)
	if err != nil {
		log.Fatal(err)
	}

	for _, linha := range C {
		fmt.Println(linha)
	}
}
