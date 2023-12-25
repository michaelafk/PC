package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/") //Conecta al servidor
	failOnError(err, "Client: Error establint la connexió amb RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel() //Obre un canal
	failOnError(err, "Client: Error al obrir un canal")
	defer ch.Close()

	deposits, err := ch.QueueDeclare( //Cua deposits
		"deposits", // name
		false,      // durable (Si el servidor cau la cua perdura amb combinació amb els missatges persistents)
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Client: Error al declarar la cua deposits")

	balances, err := ch.QueueDeclare( //Cua banlances
		"balances", // name
		false,      // durable (Si el servidor cau la cua perdura amb combinació amb els missatges persistents)
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Client : Error al declarar la cua balances")

	balMsgs, err := ch.Consume( //Missatges de la cua "balances"
		balances.Name, // queue
		"",            // consumer
		false,         // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	failOnError(err, "Client: Error al declarar els missatges de la cua balances")

	err = ch.ExchangeDeclare( // Intercanvi de tipus fanout per fer acabar als clients
		"robatori", // Nom de l'intercanvi
		"fanout",   // Tipus d'intercanvi
		true,       // Durable
		false,      // No esborrable
		false,      // Exclusiu
		false,      // No esperar a l'intercanvi ja declarat
		nil,        // Arguments addicionals
	)
	failOnError(err, "Client: Error al declarar l'intercamvi robatori")

	acabament, err := ch.QueueDeclare( // Cua temporal (esborrable quan tots els clients es tanquen)
		"",    // Nom de la cua generat automàticament
		false, // No es duradora (esborrable quan tots els consumidors es tanquen)
		true,  // Exclusiva (és eliminada quan aquest consumidor es tanca)
		false, // No esperar a la cua ja declarada
		false, // No eliminar la cua quan es tanca l'últim consumidor
		nil,   // Arguments addicionals
	)
	failOnError(err, "Client: Error al declarar cua d'acabament")

	err = ch.QueueBind( // Vincular acabament la cua a l'intercanvi robatori
		acabament.Name, // Nom de la cua
		"",             // Clau de l'intercanvi (fanout no utilitza claus)
		"robatori",     // Nom de l'intercanvi
		false,          // No esperar a la cua ja lligada
		nil,            // Arguments addicionals
	)
	failOnError(err, "Client: Error al cincular acabament a robatori")

	acaMsgs, err := ch.Consume( // Missatges de la cua acabament
		acabament.Name, // Nom de la cua
		"",             // Consumidor
		true,           // Auto-ack (marcar missatges com a llegits automàticament)
		false,          // No exclusiu
		false,          // No eliminar la cua quan es tanca l'últim consumidor
		false,          // No esperar confirmació
		nil,            // Arguments addicionals
	)
	failOnError(err, "Client: Error al declarar els missatges de la cua acabament")

	go func() { // Missatges de la cua d'acabament
		for msg := range acaMsgs {
			body := string(msg.Body)
			fmt.Println("El banquer ha dit: " + body)
			fmt.Println("Me'n vaig idò!")
			os.Exit(0)
		}
	}()

	nom := nameFrom(os.Args)
	fmt.Println("Hola el meu nom és: " + nom)
	operacions := 1 + rand.Intn(5)
	fmt.Println(nom + " vol fer " + strconv.Itoa(operacions) + " operacions")
	op := 0
	for i := 1; i <= operacions; i++ {
		signe := rand.Intn(2)
		op = rand.Intn(10) + 1 // +1 Per no fer operacions de 0
		if signe == 0 {
			op = -op
		}
		fmt.Println(nom + " operació " + strconv.Itoa(i) + ": " + strconv.Itoa(op))
		err = ch.Publish( // Pulicar el missatge a la cua dels diposits
			"",            // Nom de l'intercamvi (buit ja que publiquem directament a la cua)
			deposits.Name, // Nom de la cua
			false,         // No esperar confirmació
			false,         // No eliminar missatges no lliurats
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(fmt.Sprintf("%s;%d", nom, op)), // ; caràcter separador
			})
		failOnError(err, "Client: Error en publicar a diposits")
		fmt.Println("Operació sol·licitada")
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		for d := range balMsgs {
			total, _ := strconv.Atoi(string(d.Body))
			if total == 0 && op > 0 {
				fmt.Println("TESORER SOSPITÓS")
			} else if total == 0 && op < 0 {
				fmt.Println("NO HI HA SALDO")
			} else if total > 0 && op > 0 {
				fmt.Println("INGRÉS CORRECTE")
			} else {
				fmt.Println("ES FARÀ EL REINTEGRE SI ÉS POSSIBLE")
			}
			fmt.Println("Balanç actual: " + strconv.Itoa(total))
			d.Ack(true)
			break
		}
		fmt.Println(strconv.Itoa(i) + "----------------------------")
		time.Sleep(time.Duration(rand.Intn(9000)) * time.Millisecond)
	}
}

func nameFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "Catalina de Plaça"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}
