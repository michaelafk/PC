package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Message struct {
	IntValue int
	Name     string
}

type Feedback struct {
	Continue int
	frase    string
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Queue for sending messages
	q, err := ch.QueueDeclare(
		"compte", // name
		false,    // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare a queue")

	feedbackChannel := make(chan Feedback)

	go consumeMessages(ch, q.Name, feedbackChannel)

	nom := os.Args[1]

	log.Printf("Hola el meu nom és: %s", nom)
	operacions := rand.Intn(10) + 1
	log.Printf("%s vol fer %d operacions", nom, operacions)

	for i := 0; i < operacions; i++ {
		op := rand.Intn(40) - 20
		log.Printf("%s operació %d: %d", nom, i, op)
		message := Message{
			IntValue: op,
			Name:     nom,
		}
		sendMessage(ch, q.Name, message)
		log.Printf("Operació sol·licitada...")
		log.Printf("%d---------------------------------------", i+1)

		// Wait for feedback
		feedback := <-feedbackChannel
		if feedback.Continue == 0 {
			log.Printf("El Tesorer ha dit: %s\nJo també me'n vaig idò!", feedback.frase)
			break
		}
		time.Sleep(time.Duration(1) * time.Second)
	}
}

func sendMessage(ch *amqp.Channel, queueName string, message Message) {
	body := []byte(strconv.Itoa(message.IntValue) + "," + message.Name)

	err := ch.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	failOnError(err, "Failed to publish a message")
}

func consumeMessages(ch *amqp.Channel, queueName string, feedbackChannel chan Feedback) {
	msgs, err := ch.Consume(
		queueName, // queue
		"Client",  // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	failOnError(err, "Failed to register a consumer")

	for d := range msgs {
		parts := strings.Split(string(d.Body), ",")
		if len(parts) != 2 {
			log.Printf("Invalid message format: %s", d.Body)
			continue
		}

		intValue, err := strconv.Atoi(parts[0])
		if err != nil {
			log.Printf("Invalid integer in message: %s", parts[0])
			continue
		}

		feedback := Feedback{
			Continue: intValue,
			frase:    parts[1],
		}

		feedbackChannel <- feedback
	}
}
