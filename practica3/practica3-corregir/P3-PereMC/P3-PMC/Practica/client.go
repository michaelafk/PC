// @Author: Pere Martínez Cifre
package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

const (
	nombreOperacionsMax = 10 // Quantitat màxima d'operacions per client

	// valorMaxOperacions indica el valor de (0,n] que pot collir el generador
	// aleatori. Per poder permetre que es generin valors negatius, el valor de
	// compensadorOperacions es resta al primer de tal manera que, el rang passa
	// a ser (-n/2, n/2]. P. ex. On valorMaxOperacions = 10,
	// compensadorOperacions = 5 (10 / 2). De tal manera que si el valor
	// aleatori és 7 -> 7 - 5 = 2. Si surt 3 -> 3 - 5 = -2.
	valorMaxOperacions    = 10
	compensadorOperacions = valorMaxOperacions / 2

	// Milisegons que es vol que duri l'espera que simula el temps de l'operació
	milisegonsEspera = 50
)

// Imprimir una línia en format log a consola
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// Indiciar l'inici d'una operació
func printIntro(order int) {
	fmt.Printf("%d---------------------------\n", order+1)
}

func main() {
	// Collir el nom dels arguments extra indicats a l'execució
	// Si té un llinatge, també l'afageix al nom.
	// "go run client.go <nom> <llinatge>" -> "go run client.go Bernat Metge"
	nom := os.Args[1]
	if len(os.Args) > 2 {
		nom = nom + " " + os.Args[2]
	}

	//Quantiat d'operacions a realitzar
	nombreOperacions := rand.Intn(nombreOperacionsMax) + 1

	// Connectar amb RabbitMQ
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Error creant la conexió amb RabbitMQ")
	defer conn.Close()

	// Obrir el canal dins rabbitMQ
	ch, err := conn.Channel()
	failOnError(err, "Error obrint el canal")
	defer ch.Close()

	// Crear o connectar (Si ja s'ha creat) l'exchange necessari per
	// rebre el missatge de despedida de tipus Fanout. Cada client crearà
	// una cua única qComunicacions que també vincularà a aquest exchange.
	// D'aquesta manera, la quantitat de cues pot ser dinàmica.
	err = ch.ExchangeDeclare(
		"fanoutExchange", // nom de l'exchange
		"fanout",         // tipus d'exchange
		false,            // durable?
		true,             // auto-borrar?
		false,            // interna?
		false,            // sense-espera?
		nil,              // arguments
	)
	failOnError(err, "Error creant l'exchange fanout")

	// Crear la cua Balances, aquesta cua l'utilitzaràn els clients per indicar
	// les diferentes operacions que desitgen fer. El tesorer les cullirà d'una
	// en una per processar-les.
	qBalances, err := ch.QueueDeclare(
		"Balances", // nom de la cua
		false,      // durable?
		false,      // auto-borrar?
		false,      // exclusiva?
		false,      // sense-espera?
		nil,        // arguments
	)
	failOnError(err, "Error creant la cua Balances")

	// Crear la cua Diposits, aquesta cua l'utilitzarà el tesorer per publicar
	// el balanç del compte després de cada operació processada de Balances.
	// Els clientes culliràn els missatges i els compararàn amb l'operació
	// publicada per ells a la cua Balances.
	qDiposits, err := ch.QueueDeclare(
		"Diposits", // nom de la cua
		false,      // durable?
		false,      // auto-borrar?
		false,      // exclusiva?
		false,      // sense-espera?
		nil,        // arguments
	)
	failOnError(err, "Error creant la cua Diposits")

	// Crear la cua "Comunicacio <nom> [llinatge]". P. Ex "Comunicacio Jaume"
	// Aquesta cua s'usa exclusivament per rebre l'avís de tancament d'oficina
	nomComunicacio := "Comunicacio " + nom
	qComunicacio, err := ch.QueueDeclare(
		nomComunicacio, // nom de la cua
		false,          // durable?
		true,           // auto-borrar?
		false,          // exclusiva?
		false,          // sense-espera?
		nil,            // arguments
	)
	failOnError(err, "Error creant la cua Comunicacio")

	ch.Qos(1, 0, false) //Fixam l'enviament de missatges a un per consumidor

	// Vinculam la cua de comunicacions a l'Exchange de tal manera que es pugui
	// rebre l'avís de tancament
	err = ch.QueueBind(
		qComunicacio.Name, // nom de la cua
		"",                // tipus de routing
		"fanoutExchange",  // exchange
		false,             // sense-espera?
		nil,               // arguments
	)
	failOnError(err, "Error vinculat la cua de comunicacions amb l'exchange")

	// Crear relació de tipus consumidor amb la cua Dipòsits. Això permet al
	// client collir les respostes a les seves operacions amb el nou balanç
	missatgesDiposits, err := ch.Consume(
		qDiposits.Name, // nom de la cua
		"",             // identificador del consumidor
		true,           // auto-ack?
		false,          // exclusiva?
		false,          // no-local?
		false,          // sense-espera?
		nil,            // args
	)
	failOnError(err, "Error registrant el consumidor per Diposits")

	// Crear relació de tipus consumidor amb la cua pròpia de Comunicacions.
	// Això permet al client collir l'avís quan s'ha tancat l'oficina
	missatgeComunicacio, err := ch.Consume(
		qComunicacio.Name, // nom de la cua
		"",                // identificador del consumidor
		true,              // auto-ack?
		false,             // exclusiva?
		false,             // no-local?
		false,             // sense-espera?
		nil,               // args
	)
	failOnError(err, "Error registrant el consumidor per Comunicacions")

	// El client es presenta
	fmt.Printf("Hola el meu nom és : %s\n", nom)
	fmt.Printf("%s vol fer %d operacions\n", nom, nombreOperacions)

	// Per cada operació que vulgui fer
	for i := 0; i < nombreOperacions; i++ {
		printIntro(i)
		// Agafa un valor aleatori (-n/2, n/2] que no sigui 0.
		operacioI := 0
		for {
			operacioI = rand.Intn(valorMaxOperacions) - compensadorOperacions

			if operacioI != 0 {
				break
			}
		}

		// Indica l'operació
		fmt.Printf("%s operació %d: %d\n", nom, i+1, operacioI)

		// Publica a la cua Balances el missatge amb format
		// "valorOperacio|nomClient"
		body := strconv.Itoa(operacioI) + "|" + nom
		ctx := context.TODO()
		err = ch.PublishWithContext(
			ctx,            // contexte (genèric en el nostre cas)
			"",             // exchange (predeterminat per arribar a una cua)
			qBalances.Name, // cua objectiu
			false,          // isMandatory?
			false,          // isImmediate?
			amqp091.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		failOnError(err, "Error publicant a la cua Balances")

		// Indica que ha sol·licitat l'operació
		fmt.Printf("Operació sol·licitada...\n")

		// Espera n milisegons
		time.Sleep(time.Millisecond * milisegonsEspera)

		// Cull la resposta a la cua depòsits
		respostaDiposits := <-missatgesDiposits

		// Si la resposta es un string buit, vol dir que la cua ja ha estat borrada
		if string(respostaDiposits.Body) == "" {

			// Revisa que s'hagi rebut l'avís
			respostaComunicacio := string((<-missatgeComunicacio).Body)
			fmt.Printf("El banquer ha dit: %s\n", respostaComunicacio)
			fmt.Printf("Me'n vaig idò!\n\n")
			break
		}

		// Convertir la resposta amb el valor del balanç a un enter
		balanç, err := strconv.Atoi(string(respostaDiposits.Body))
		failOnError(err, "Error convertint a enter")

		// Indica el resultat de l'operació
		if operacioI > 0 && balanç == 0 {
			fmt.Printf("EL TESORER SEMBLA SOSPITÓS\n")
		} else if balanç == 0 {
			fmt.Printf("NO HI HA SALDO AL COMPTE\n")
		} else if operacioI > 0 {
			fmt.Printf("INGRÉS CORRECTE\n")
		} else {
			fmt.Printf("ES FARÀ EL REINTEGRE SI HI HA SALDO\n")
		}

		// S'indica el balanç que s'ha rebut a la resposta
		fmt.Printf("Balanç actual: %d\n", balanç)
	}

}
