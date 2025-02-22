package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	matrixSize = 10
)

var (
	matrixA = [matrixSize][matrixSize] int {}
	matrixB = [matrixSize][matrixSize] int {}
	result = [matrixSize][matrixSize] int{}
	rwLock = sync.RWMutex{}
	cond = sync.NewCond(rwLock.RLocker())
	waitGroup = sync.WaitGroup{}
)

func generateRandomMatrix(matrix *[matrixSize][matrixSize] int) {
	for row := 0; row < matrixSize; row++ {
		for col := 0; col < matrixSize; col++ {
			matrix[row][col] += rand.Intn(10) - 5
		}
	}
}

func workOutRow(row int) {
	rwLock.RLocker()
	for {
		waitGroup.Done()
		cond.Wait()
		for col := 0; col < matrixSize; col++ {
			for i := 0; i < matrixSize; i++ {
				result[row][col] += matrixA[row][i] * matrixB[i][col]
			}
		}
	}

}

func main() {
	fmt.Println("Working...")
	waitGroup.Add(matrixSize)
	for row := 0; row < matrixSize; row++ {
		go workOutRow(row)
		//fmt.Printf("%v\n", result[row])
	}
	start := time.Now()
	for i := 0; i < 100; i++ {
		waitGroup.Wait()
		rwLock.Lock()
		generateRandomMatrix(&matrixA)
		generateRandomMatrix(&matrixB)
		waitGroup.Add(matrixSize)
		rwLock.Unlock()
		cond.Broadcast()
	}
	elapsed := time.Since(start)
	fmt.Printf("Done! Processing took %s", elapsed)
}


