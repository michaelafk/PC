package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	pSemaphore := make(chan int, 2)
	qSemaphore := make(chan int, 2)
	rSemaphore := make(chan int, 0)

	wg.Add(3)

	go func() {
		defer wg.Done()
		for i := 0; i < 12; i++ {
			<-pSemaphore
			fmt.Print("p\n")
			rSemaphore <- 1
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 12; i++ {
			<-qSemaphore
			fmt.Print("q\n")
			rSemaphore <- 1
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 20; i++ {
			<-rSemaphore
			fmt.Print("r\n")
			if i%2 == 0 {
				pSemaphore <- 1
			} else {
				qSemaphore <- 1
			}
		}
	}()

	pSemaphore <- 1
	qSemaphore <- 1
	wg.Wait()
}
