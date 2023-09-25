package main

import "fmt"

func main() {
	var myarray [3]int

	for i := 0; i < len(myarray); i++ {
		myarray[i] = i + 1
	}
	for j := 0; j < len(myarray); j++ {
		fmt.Println(myarray[j])
	}
}
