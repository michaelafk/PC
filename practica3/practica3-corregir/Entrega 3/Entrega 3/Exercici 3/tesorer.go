package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

type Operacion struct {
	Nombre string
	Valor  int
}

func main() {

	var n, balance int
	var nombre string

	// Configuración de conexión con RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "No se pudo conectar a RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "No se pudo abrir el canal")
	defer ch.Close()

	// Crear cola de depósitos y balances
	depositQueue2, err := ch.QueueDeclare(
		"Deposits2", // Nombre de la cola
		false,       // Durabilidad
		false,       // Autoeliminación
		false,       // Exclusividad
		false,       // No espera para publicar
		nil,         // Argumentos adicionales
	)
	failOnError(err, "No se pudo crear la cola de depósitos")

	balanceQueue2, err := ch.QueueDeclare(
		"Balances2", // Nombre de la cola
		false,       // Durabilidad
		false,       // Autoeliminación
		false,       // Exclusividad
		false,       // No espera para publicar
		nil,         // Argumentos adicionales
	)
	failOnError(err, "No se pudo crear la cola de balances")

	exchangeName := "fanout_exchange"
	err = ch.ExchangeDeclare(
		exchangeName, "fanout", true, false, false, false, nil,
	)
	failOnError(err, "Failed to declare the fanout exchange")

	minBoti := 20 // Mínimo botín esperado
	balance = 0   // Valor inicial del balance

	fmt.Printf("El tresorer és al despatx. El botí mínim és: %d\n", minBoti)
	// Consumir mensajes de la cola de depósitos
	msgs, err := ch.Consume(
		depositQueue2.Name, // Nombre de la cola
		"",                 // Nombre del consumidor
		false,              // Auto-Ack
		false,              // Exclusividad
		false,              // No espera para consumir
		false,              // Argumentos adicionales
		nil,
	)
	failOnError(err, "No se pudo consumir la cola de depósitos")

	// Manejar mensajes de depósitos
	for msg := range msgs {

		fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

		var operacion Operacion
		err := json.Unmarshal(msg.Body, &operacion)
		failOnError(err, "Failed to unmarshal JSON")

		n = operacion.Valor
		nombre = operacion.Nombre

		fmt.Printf("Operació rebuda: %d del client: %s\n", n, nombre)

		msg.Ack(false)

		if n >= 0 {
			balance += n
		} else if balance+n >= 0 {
			balance += n
		} else {
			fmt.Println("OPERACIÓ NO PERMESA NO HI HA FONS")
		}

		fmt.Printf("Balanç: %d\n", balance)

		if balance >= minBoti {
			fmt.Println("El Tresorer decideix robar el diposit i tancar el despatx")

			// Enviar mensaje a todos los clientes activos
			err := ch.Publish(
				"",
				balanceQueue2.Name, // Envía a todos los enlaces
				false,              // Mandar como persistente
				false,              // No esperar confirmación
				amqp.Publishing{
					DeliveryMode: amqp.Persistent,
					ContentType:  "text/plain",
					Body:         []byte(strconv.Itoa(0)),
				},
			)
			failOnError(err, "No se pudo enviar el mensaje de cierre a los clientes")

			fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
			fmt.Println("El Tresorer se'n va\n\n")

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			bodyF := "L'Oficina acaba de tancar\n"

			err = ch.PublishWithContext(ctx,
				exchangeName,
				"",
				false,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(bodyF),
				})
			failOnError(err, "Failed to publish a message")

			// Eliminar todas las colas del sistema
			ch.QueueDelete(depositQueue2.Name, false, false, false)
			ch.QueueDelete(balanceQueue2.Name, false, false, false)
			fmt.Printf("Cua '%s' esborrada amb èxit.\n", depositQueue2.Name)
			fmt.Printf("Cua '%s' esborrada amb èxit.\n", balanceQueue2.Name)

			os.Exit(0)
		} else {
			// Enviar mensaje a todos los clientes activos
			err := ch.Publish(
				"",
				balanceQueue2.Name, // Envía a todos los enlaces
				false,              // Mandar como persistente
				false,              // No esperar confirmación
				amqp.Publishing{
					DeliveryMode: amqp.Persistent,
					ContentType:  "text/plain",
					Body:         []byte(strconv.Itoa(balance)),
				},
			)
			failOnError(err, "No se pudo enviar el mensaje de cierre a los clientes")

			fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")

		}

	}
}
