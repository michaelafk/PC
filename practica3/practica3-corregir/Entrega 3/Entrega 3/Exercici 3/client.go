package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
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

	args := os.Args

	if len(args) > 1 {
		// Configuración de conexión con RabbitMQ
		conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
		failOnError(err, "No se pudo conectar a RabbitMQ")
		defer conn.Close()

		ch, err := conn.Channel()
		failOnError(err, "No se pudo abrir el canal")
		defer ch.Close()

		// Mensaje de inicio del cliente
		nombre := args[1] // Nombre del cliente, puedes cambiarlo si es necesario
		rand.Seed(time.Now().UnixNano())
		operaciones := rand.Intn(5) + 1 // Cantidad de operaciones a realizar por el cliente

		fmt.Printf("Hola, el meu nom és: %s\n", nombre)
		fmt.Printf("%s vol fer %d operacions\n", nombre, operaciones)

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

		q, err := ch.QueueDeclare(
			"",    //nombre
			false, //durable
			false, //delete when unused
			true,  //exclusive
			false, //no-wait
			nil,   //arguments
		)
		failOnError(err, "Failed to declare a queue")

		err = ch.QueueBind(
			q.Name,       //nombre de la cola
			"",           //routing key
			exchangeName, //exchange
			false,
			nil,
		)
		failOnError(err, "Failed to bind a queue")

		msgs, err := ch.Consume(
			q.Name, //queue
			"",     //consumer
			false,  //auto ack
			false,  //exclusive
			false,  //no local
			false,  //no wait
			nil,    //args
		)
		failOnError(err, "Failed to register a consumer")

		go func() {
			for {
				d := <-msgs
				if string(d.Body) == "L'Oficina acaba de tancar\n" {
					fmt.Printf("El banquer ha dit: %s", string(d.Body))
					fmt.Printf("Jo també m'en vaig idò!\n")
					os.Exit(0)
				}
			}
		}()

		for i := 1; i <= operaciones; i++ {

			time.Sleep(2 * time.Second)
			tipo := rand.Intn(2)      //genera numero aleatorio entre 0 (retiro) y 1 (ingreso)
			cantidad := rand.Intn(20) //generar numero aleatorio entre 0 y 10

			if tipo == 0 {
				cantidad = cantidad * (-1)
			}

			fmt.Printf("%s operació %d: %d\n", nombre, i, cantidad)

			operacion := Operacion{nombre, cantidad}

			body, err := json.Marshal(operacion)
			failOnError(err, "Failed to Marshal operation")

			err = ch.Publish(
				"",
				depositQueue2.Name,
				false,
				false,
				amqp.Publishing{
					DeliveryMode: amqp.Persistent,
					ContentType:  "aplication/json",
					Body:         body,
				})
			failOnError(err, "Failed to publish a message to the first queue")

			fmt.Printf("Operació solicitada...\n")
			time.Sleep(5 * time.Second)
			//esperar balance
			msg, ok, err := ch.Get(balanceQueue2.Name, false)
			failOnError(err, "Failed to get message from queue")

			if ok {
				balance, _ := strconv.Atoi(string(msg.Body))

				if tipo == 1 && balance == 0 {
					fmt.Printf("EL TESORER SEMBLA SOSPITÓS\n")
				} else if balance == 0 {
					fmt.Printf("NO HI HA SALDO AL COMPTE\n")
				} else if tipo == 1 {
					fmt.Printf("INGRÉS CORRECTE\n")
				} else {
					fmt.Printf("ES FARÀ EL REINTEGRE SI HI HA SALDO\n")
				}

				fmt.Printf("Balanç actual: %d\n", balance)

				err = msg.Ack(false)
				failOnError(err, "Failed to acknowlege the message")
			}

			fmt.Printf("%d-------------------------\n", i)
		}
	} else {
		fmt.Printf("Proporcione un nombre como argumento.\n")
	}

}
