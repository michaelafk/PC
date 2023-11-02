package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	Abelles := make(chan int, 1)
	Os := make(chan int, 1)
}
