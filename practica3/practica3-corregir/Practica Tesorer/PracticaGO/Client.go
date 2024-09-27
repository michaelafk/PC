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
	amqp "github.com/rabbitmq/amqp091-go"
)

// Variables
var Nom string
var NombreOperacions int = 1
var MissatgeOperacio string = ""
var EsIngres bool // variable per saber si s'ha realitzar un ingrés o no
var BalançTresorer int

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func fer_operació(conn *amqp091.Connection) {

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	corrId := randomString(32)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",         // exchange
		"Diposits", // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          []byte(MissatgeOperacio),
		})
	failOnError(err, "Failed to publish a message")

	fmt.Println("Operació sol·licitada...")
	time.Sleep(2 * time.Second)

	//Rebuda de respostes
	for d := range msgs {
		if corrId == d.CorrelationId {
			BalançTresorer, err = strconv.Atoi(string(d.Body))
			if BalançTresorer == 0 {
				if EsIngres {
					fmt.Println("EL TESORER SEMBLA SOSPITÓS")
					fmt.Println("El Tesorer ha dit: L'oficina acaba de tancar")
					fmt.Println("Jo també me'n vaig idò")
					os.Exit(0)

				} else {
					fmt.Println("NO HI HA SALDO AL COMPTE")
				}

			} else {
				if EsIngres {
					fmt.Println("INGRÉS CORRECTE")
				} else {
					fmt.Println("ES FARÀ EL REINTEGRE SI HI HA SALDO")
				}
			}

			fmt.Println("Balanç actual: " + string(d.Body))
			break
		}
	}

	return
}

func main() {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	//Lectura del nom per consola
	Nom = os.Args[1]
	//Missatges inici de simulació
	fmt.Println("Hola el meu nom és: " + Nom)
	rand.Seed(time.Now().UnixNano())
	//Generam el nombre aleatori
	NombreOperacions := rand.Intn(5) + 1
	fmt.Println(Nom + " vol fer " + strconv.Itoa(NombreOperacions) + " operacions")
	for i := 0; i < NombreOperacions; i++ {
		generar_missatge_diposit(i)
		fer_operació(conn)
		rebreMissatge(conn)
	}
}

func generar_missatge_diposit(i int) {
	rand.Seed(time.Now().UnixNano())
	//Generam el nombre aleatori
	ValorOperacio := rand.Intn(9) + 1
	//Convertim el nombre en String
	ValorOperacioStr := strconv.Itoa(ValorOperacio)
	//Generam el missatge Operacio amb un 50% de possibilitats de que sigui un ingrés o un retir de diners
	randomAux := rand.Float64()
	if randomAux < 0.5 {
		//La coma serveix de separador
		MissatgeOperacio = Nom + ",+" + ValorOperacioStr
		EsIngres = true
		fmt.Println(Nom + " operació " + strconv.Itoa(i+1) + ": +" + ValorOperacioStr)
	} else {
		MissatgeOperacio = Nom + ",-" + ValorOperacioStr
		EsIngres = false
		fmt.Println(Nom + " operació " + strconv.Itoa(i+1) + ": -" + ValorOperacioStr)
	}

}

func rebreMissatge(conn *amqp091.Connection) {

	canalFanout, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer canalFanout.Close()

	err = canalFanout.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := canalFanout.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = canalFanout.QueueBind(
		q.Name, // queue name
		"",     // routing key
		"logs", // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	msgs, err := canalFanout.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			fmt.Println(string(d.Body))
		}
	}()

}
