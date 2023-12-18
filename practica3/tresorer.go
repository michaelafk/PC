package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	MininBoti = 20
	Separador = "$"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func Operacion(n int, total *int) int {
	if n > 0 { //ingreso
		*total = *total + n
		return 1
	} else if n < 0 { //reintegro
		if (*total + n) >= 0 { //reintegro valido
			*total = *total + n
			return 2
		} else {
			return -1 //reintegro no valido
		}
	}
	return 0
}

func split(r rune) bool {
	return r == '$'
}

func main() {
	total := 0
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"Deposits", // name
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
			//n, err := strconv.Atoi(string(d.Body))
			//failOnError(err, "Failed to convert body to integer")
			compositeMsg := strings.FieldsFunc(string(d.Body), split)
			usuario := compositeMsg[0]
			cantidad, err := strconv.Atoi(compositeMsg[1])
			//str_concat, err := strings.Join(str_slices, "-");
			failOnError(err, "Failed to convert compositeMsg[1]")
			fmt.Printf(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>\n")
			fmt.Printf("Operacio rebuda: %d del client: %s\n", cantidad, usuario)
			estado := Operacion(cantidad, &total)
			var str strings.Builder
			switch estado {
			case 1:
				//INGRES CORRECTE
				str.WriteString("INGRES CORRECTE")
			case 2:
				//ES FARA EL REINTEGRE SI ES POSSIBLE
				str.WriteString("ES FARA EL REINTEGRE SI ES POSSIBLE")
			case -1:
				//NO HI HA SALDO
				str.WriteString("NO HI HA SALDO")
				fmt.Println("OPERACIO NO PERMESA NO HI HA FONS")
			}
			fmt.Printf("Balanç: %d\n", total)
			fmt.Printf("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<\n")
			str.WriteString(Separador)
			str.WriteString(strconv.Itoa(total))
			err = ch.PublishWithContext(ctx,
				"",        // exchange
				d.ReplyTo, // routing key
				false,     // mandatory
				false,     // immediate
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          []byte(str.String()),
				})
			failOnError(err, "Failed to publish a message")

			d.Ack(false)
		}
	}()

	log.Printf(" [*] Esperant clients")
	<-forever
}
