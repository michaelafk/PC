package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Variables
const MinimBoti int = 20

var Balanç int = 0

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {

	//Missatge inicial
	fmt.Println("El tresorer és al despatx. El botí mínim és: " + strconv.Itoa(MinimBoti))

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"Diposits", // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		for d := range msgs {
			//variable per manejar els missatges rebuts
			var missatgeSeparat []string
			//Passam el missatge a String que te forma de "Nom,SigneValor" on el "signe"
			//pot ser "+" o "-" amb un 50% de porobabilitats
			//I separam el missatge per el separador "," y el guardam
			missatgeSeparat = strings.Split(string(d.Body), ",")

			//Manejam la lògica del missatge
			manejar_missatge(missatgeSeparat[0], missatgeSeparat[1])

			//Cas de que s'arribi al mínim botí
			if Balanç >= MinimBoti {

				fmt.Println("El Tresorer decideix robar el diposit i tancar el despatx")
				Balanç = 0

				canalFanout, err := conn.Channel()
				failOnError(err, "Failed to open a channel")
				defer canalFanout.Close()

				err = canalFanout.ExchangeDeclare(
					"logs",   // name
					"fanout", //type
					true,
					false,
					false,
					false,
					nil,
				)
				failOnError(err, "Failed to declare exchange")
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				err = canalFanout.PublishWithContext(ctx,
					"logs", // exchange
					"",     // routing key
					false,  // mandatory
					false,  // immediate
					amqp.Publishing{
						ContentType: "text/plain",
						Body:        []byte("L'oficina acaba de tancar"),
					})
				failOnError(err, "Failed to publish a message")

				fmt.Println("El Tresorer se'n va\n\n")

				//La resposta és el balanç amb format string
				response := strconv.Itoa(Balanç)

				err = ch.PublishWithContext(ctx,
					"",        // exchange
					d.ReplyTo, // routing key
					false,     // mandatory
					false,     // immediate
					amqp.Publishing{
						ContentType:   "text/plain",
						CorrelationId: d.CorrelationId,
						Body:          []byte(response),
					})
				failOnError(err, "Failed to publish a message")

				d.Ack(false)

				os.Exit(0)
			}

			//La resposta és el balanç amb format string
			response := strconv.Itoa(Balanç)

			err = ch.PublishWithContext(ctx,
				"",        // exchange
				d.ReplyTo, // routing key
				false,     // mandatory
				false,     // immediate
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          []byte(response),
				})
			failOnError(err, "Failed to publish a message")

			d.Ack(false)

		}
	}()

	log.Printf(" [*] Awaiting RPC requests")
	<-forever
}

func manejar_missatge(nom string, operacio string) {
	//Print sobre el missatge rebut
	fmt.Println("Operacio rebuda: " + operacio + " del client: " + nom)

	val, err := strconv.Atoi(string(operacio[1]))
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	//Lògica sobre incrementar o decrementar el balanç
	if operacio[0] == '+' {
		Balanç += val

	} else {
		if val > Balanç {
			fmt.Println("OPERACIÓ NO PERMESA, NO HI HA FONS")
		} else {
			Balanç -= val
		}
	}

	fmt.Println("Balanç: " + strconv.Itoa(Balanç))
}
