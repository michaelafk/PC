package main

import (
    "fmt"
    "sync"
)

var (
    turn  int
    flag0 bool
    flag1 bool
)

func process0() {
    for {
        flag0 = true
        turn = 1
        for flag1 && turn == 1 {
            // Espera activa
        }
        // Sección crítica para el proceso 0
        fmt.Println("Proceso 0 en la sección crítica")
        flag0 = false
        // Sección no crítica
    }
}

func process1() {
    for {
        flag1 = true
        turn = 0
        for flag0 && turn == 0 {
            // Espera activa
        }
        // Sección crítica para el proceso 1
        fmt.Println("Proceso 1 en la sección crítica")
        flag1 = false
        // Sección no crítica
    }
}

func main() {
    flag0 = false
    flag1 = false
    turn = 0

    var wg sync.WaitGroup
    wg.Add(2)

    go func() {
        defer wg.Done()
        process0()
    }()

    go func() {
        defer wg.Done()
        process1()
    }()

    wg.Wait()
}
