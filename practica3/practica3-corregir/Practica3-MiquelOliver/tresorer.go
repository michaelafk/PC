package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"strconv"
	"strings"
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
	MinimBoti := 50
	balance := 0

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"compte", // name
		false,    // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare a queue")

	messageChannel := make(chan Message)

	go consumeMessages(ch, q.Name, messageChannel)
	log.Printf("El tresorer és al despatx. El botí mínim és: " + strconv.Itoa(MinimBoti))

	var forever chan struct{}

	go func() {
		for msg := range messageChannel {
			log.Printf(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
			log.Printf("Operació rebuda: %d del client: %s", msg.IntValue, msg.Name)

			feedback := Feedback{1, ""}

			if msg.IntValue >= 0 {
				balance += msg.IntValue
				log.Printf("INGRÉS CORRECTE")
				log.Printf("Balanç actual: %d", balance)
				if balance == 0 {
					log.Printf("EL TESORER SEMBLA SOSPITÓS")
				}
			} else {
				if balance == 0 {
					log.Printf("NO HI HA SALDO AL COMPTE")
				} else {
					if balance+msg.IntValue < 0 {
						log.Printf("OPERACIÓ NO PERMESA NO HI HA FONS")
					} else {
						log.Printf("ES FARÀ EL REINTEGRE SI HI HA SALDO")
						balance += msg.IntValue
					}
				}
				log.Printf("Balanç actual: %d", balance)
			}

			if balance >= MinimBoti {
				feedback.Continue = 0
				feedback.frase = "L'oficina acaba de tancar"
				log.Printf("El Tresorer decideix robar el diposit i tancar el despatx")
				log.Printf("El Tresorer diu \"%s\" i se'n va\n\n", feedback.frase)
				sendMessage(ch, q.Name, feedback)

				if _, err := ch.QueueDelete(q.Name, false, false, false); err != nil {
					log.Printf("Failed to delete queue %s: %s\n", q.Name, err)
				} else {
					log.Printf("Cua '%s' esborrada amb èxit.\n", q.Name)
				}
				break
			}
			sendMessage(ch, q.Name, feedback)
			log.Printf("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
		}
	}()

	if balance < MinimBoti {
		log.Printf(" [*] Esperant clients")
		<-forever
	}
}

func sendMessage(ch *amqp.Channel, queueName string, message Feedback) {
	body := []byte(strconv.Itoa(message.Continue) + "," + message.frase)

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

func consumeMessages(ch *amqp.Channel, queueName string, messageChannel chan Message) {
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

		message := Message{
			IntValue: intValue,
			Name:     parts[1],
		}

		messageChannel <- message
	}
}
