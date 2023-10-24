package main

import (
	"fmt"
	"sync"
	"time"
)

var flag = []bool{false, false}
var last int
var wg sync.WaitGroup
var Entrades [2]int
var entradesTotals int = 5

func Entrada(id int) {
	defer wg.Done()

	//hacer algo
	fmt.Printf("Entrada %d\n", id)
	Entrades[id] = 0
	for i := 0; i < entradesTotals; i++ {
		time.Sleep(time.Millisecond * 500)
		flag[id] = true
		last = id
		var pos int = (id + 1) % 2
		for flag[pos] && last == id {
			//espera activa
		}
		//seccion critica
		fmt.Printf("Porta %d: ", id)
		Entrades[id] += 1
		fmt.Printf("%d entrades de : %d Temps: %T", Entrades[id], i, time.Now())
		fmt.Printf("Porta %d\n", id)

		flag[id] = false

		time.Sleep(time.Millisecond * 200)
	}
}
