package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	Separador = "$"
)

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

func split(r rune) bool {
	return r == '$'
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func OperacionRPC(Nombre string, cantidad int) (Estado string, total int) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

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
	//rand.Seed(time.Now().UTC().UnixNano())
	var str strings.Builder
	str.WriteString(Nombre)
	str.WriteString(Separador)
	str.WriteString(strconv.Itoa(cantidad))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",         // exchange
		"Deposits", // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          []byte(str.String()),
		})
	failOnError(err, "Failed to publish a message")

	for d := range msgs {
		if corrId == d.CorrelationId {
			//res, err = strconv.Atoi(string(d.Body))
			compositeMsg := strings.FieldsFunc(string(d.Body), split)
			Estado = compositeMsg[0]
			total, err = strconv.Atoi(compositeMsg[1])

			failOnError(err, "Failed to convert compositeMsg[1] to integer")
			break
		}
	}

	return
}

func main() {
	//rand.Seed(time.Now().UTC().UnixNano())

	Nombre := os.Args[1]
	NumeroOperaciones := randInt(1, 5)
	fmt.Printf("Hola el meu nom es: %s\n", Nombre)
	fmt.Printf("%s vol fer %d operacions\n", Nombre, NumeroOperaciones)
	for i := 0; i < NumeroOperaciones; i++ {
		cantidad := randInt(-7, 7)
		fmt.Printf("%s operacio %d: %d\n", Nombre, i+1, cantidad)
		fmt.Printf("Operacio sol-licitada\n")
		Estado, total := OperacionRPC(Nombre, cantidad)
		fmt.Printf("%s\n", Estado)
		fmt.Printf("Balanç actual: %d\n", total)
		fmt.Printf("%d-------------------------\n", i+1)
	}
}
