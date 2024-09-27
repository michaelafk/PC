package main

import (
	"fmt"
	//"runtime"
	"strconv"
	"strings"
)

type Empty struct{}

const (
	NombreSalons             = 3
	NombreTaules             = 3
	NombreProcesosFumadors   = 6
	NombreProcesosNoFumadors = 6
	BufferSize               = 10
)

func Fumador(id string, done chan Empty, peticioEntrada chan string, confirmacio chan string) {
	fmt.Printf("Hola, el meu nom és %s, Voldria dinar i no vull fums\n", id)
	data := fmt.Sprintf("%s/Fum", id)
	peticioEntrada <- data
	data2 := <-confirmacio
	i, err := strconv.Atoi(data2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s diu: M'agrada molt el salo %d\n", id, i)
	done <- Empty{}
}
func NoFumador(id string, done chan Empty, peticioEntrada chan string, confirmacio chan string) {
	fmt.Printf("Hola, el meu nom és %s, voldria dinar i fumar\n", id)
	data := fmt.Sprintf("%s/NoFum", id)
	peticioEntrada <- data
	data2 := <-confirmacio
	i, err := strconv.Atoi(data2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s diu: M'agrada molt el salo %d\n", id, i)
	done <- Empty{}
}
func Meitre(done chan Empty, peticioEntrada chan string, confirmacio chan string) {
	fmt.Printf("Maitre arriba al restaurant\n")
	salons := [NombreSalons]int{0, 0, 0}
	tipus := [NombreSalons]string{"", "", ""}
	//esperant := []string{}
	for i := 0; i < NombreProcesosFumadors+NombreProcesosNoFumadors; i++ {
		data := <-peticioEntrada
		tokens := strings.Split(data, "/")
		aux := ""
		if strings.Compare(tokens[1], "Fum") == 0 {
			aux = "Fumador"
		} else {
			aux = "No Fumador"
		}
		fmt.Printf("*****Peticio rebuda de: %s que es %s\n", tokens[0], aux)
		data1 := ""
		for j := 0; j < len(salons); j++ {
			if salons[j] == 0 {
				data1 = fmt.Sprintf("%d", j)
				tipus[j] = aux
				salons[j] += 1
				break
			} else if (salons[j] < NombreSalons) && (strings.Compare(tipus[j], aux) == 0) {
				data1 = fmt.Sprintf("%d", j)
				salons[j] += 1
				break
			}
		}
		p, err := strconv.Atoi(data1)
		if err != nil {
			panic(err)
		}
		if p != -1 {
			fmt.Printf("*****El comensal %s te taula al salo %s que es de %s\n", tokens[0], data1, tipus[p])
		} else {
			fmt.Printf("*****El comensal %s no te taula degut a plenitud o tipus diferent de salo\n", tokens[0])
		}
		confirmacio <- data1
	}
	done <- Empty{}
}

func main() {
	fmt.Print("SIMULACIO DEL RESTAURANT DELS SOLITARIS\n")
	done := make(chan Empty, 1)
	peticioEntrada := make(chan string, BufferSize)
	confirmacio := make(chan string, BufferSize)
	Noms := [12]string{"Noelia", "Fermi", "Fons", "Narcisa", "Nataxa",
		"Francesc", "Nadia", "Franc", "Nuria", "Feliç", "Neus", "Facund"}
	go Meitre(done, peticioEntrada, confirmacio)
	for i := 0; i < NombreProcesosFumadors+NombreProcesosNoFumadors; i++ {
		if i%2 == 0 {
			go Fumador(Noms[i], done, peticioEntrada, confirmacio)
		} else {
			go NoFumador(Noms[i], done, peticioEntrada, confirmacio)
		}
	}
	for i := 0; i < NombreProcesosFumadors+NombreProcesosNoFumadors+1; i++ {
		<-done
	}
	fmt.Print("SIMULACIO ACABADA\n")
}
