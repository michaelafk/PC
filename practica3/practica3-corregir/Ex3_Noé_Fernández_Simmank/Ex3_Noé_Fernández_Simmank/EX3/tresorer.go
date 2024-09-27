// Noé Fernández Simmank
package main

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	BotinMinimo = 15
)

type Operacion struct {
	NombreCliente string
	Cantidad      int
}

type Empty struct{}

func tresorer(ch *amqp.Channel, ch2 *amqp.Channel, colaDiposits amqp.Queue, colaBalançes amqp.Queue, done chan Empty) {

	var botin int = 0

	fmt.Println(fmt.Sprintf("El tresorer és al despatx. El botí mínim és: %d", BotinMinimo))
	//se abre el consumo para recibir los mensages de depositos
	msgD, err := ch.Consume(
		colaDiposits.Name, // name
		"espera",          // id
		false,
		false,
		false,
		false,
		nil,
	)
	checkFail(err, "Error al abrir cola de espera")

	//por cada mensage que recibe
	for msg := range msgD {

		// Analizar el mensaje para obtener el nombre del cliente y la cantidad
		operacion, err := parseOperation(msg.Body)
		if err != nil {
			checkFail(err, "Error al parsear operacion")
			continue
		}
		fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
		fmt.Println(fmt.Sprintf("Operació rebuda: %d del client: %s", operacion.Cantidad, operacion.NombreCliente))
		//se hace la transaccion si es posible
		if operacion.Cantidad < 0 {
			var balanceTemp = botin + operacion.Cantidad
			if balanceTemp < 0 {
				fmt.Println("OPERACIÓ NO PERMESA NO HI HA FONS")
			} else {
				botin = botin + operacion.Cantidad
			}
		} else {
			botin = botin + operacion.Cantidad
		}
		fmt.Println(fmt.Sprintf("Balanç: %d", botin))

		//publica el saldo en la cola de balanzes
		//si es robo inidca que el saldo es 0
		if botin >= BotinMinimo {
			err2 := ch.Publish(
				"",
				colaBalançes.Name,
				false,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(strconv.Itoa(0)),
				})
			checkFail(err2, "Error al publicar un mensaje")
		} else {
			//si no publica el saldo correcto
			err2 := ch.Publish(
				"",
				colaBalançes.Name,
				false,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(strconv.Itoa(botin)),
				})
			checkFail(err2, "Error al publicar un mensaje")
		}
		fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
		if botin >= BotinMinimo {
			fmt.Println("El tresorer decideix robar el deposit i tancar el despatx")
			// Enviar mensaje de finalización a los clientes
			// duración del mensaje

			err := ch2.Publish(
				"espera",
				"",
				false,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte("L'oficina acaba de tancar"),
				})
			checkErr(err)
			fmt.Println("El Tresorer se'n va\n\n")
			break
		}
	}

	time.Sleep(5 * time.Second)

	_, err2 := ch.QueueDelete(
		colaBalançes.Name, // queue
		false,             // ifUnused
		false,             // ifEmpty
		false,             // noWait
	)
	checkFail(err2, "Error al borrar colaDiposits")
	fmt.Println(fmt.Sprintf("Cua '%s' esborrada amb èxit", colaDiposits.Name))

	_, err3 := ch.QueueDelete(
		colaDiposits.Name, // queue
		false,             // ifUnused
		false,             // ifEmpty
		false,             // noWait
	)
	checkFail(err3, "Error al borrar colaBalanç")
	fmt.Println(fmt.Sprintf("Cua '%s' esborrada amb èxit", colaBalançes.Name))

	done <- Empty{}
}

// separa el valor que recibe en body a un struct con dos elementos
func parseOperation(body []byte) (Operacion, error) {
	parts := bytes.SplitN(body, []byte(", "), 2)
	if len(parts) != 2 {
		return Operacion{}, fmt.Errorf("")
	}
	cant, err := strconv.Atoi(string(parts[0]))
	if err != nil {
		return Operacion{}, fmt.Errorf("invalid amount: %v", err)
	}

	return Operacion{
		NombreCliente: string(parts[1]),
		Cantidad:      cant,
	}, nil
}

func main() {

	// canal proceso (mensajes vacíos)
	done := make(chan Empty, 1)

	// conectarse al servidor de RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	//conn, err := amqp.Dial(rabbitMQURL)
	checkFail(err, "Error al conectarse a RabbitMQ")
	defer conn.Close()

	// abrir un canal para la cola de depositos
	ch, err := conn.Channel()
	checkFail(err, "Error al abrir un canal")
	defer ch.Close()

	// declara la cola de depositos
	colaDiposits, err := ch.QueueDeclare(
		"colaDiposits", // nombre de la cola
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            //arguments
	)
	checkFail(err, "Error al crear una cola")

	// declara la cola de balances
	colaBalançes, err := ch.QueueDeclare(
		"colaBalanç", // nombre de la cola
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          //arguments
	)
	checkFail(err, "Error al crear una cola")

	// abrir un canal tercer canal para cuando se va el espera
	ch2, err := conn.Channel()
	checkFail(err, "Error al crear un canal")
	defer ch2.Close()

	// fanout exchange -> todos los mesajes que envia el espera le
	// llegan a todos los clientes
	err = ch2.ExchangeDeclare(
		"espera",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)

	go tresorer(ch, ch2, colaDiposits, colaBalançes, done)

	// espera a que mande el mensaje por el canal indicando que ha finalizado para terminar el programa
	<-done
}

func checkFail(err error, msg string) {
	if err != nil {
		fmt.Println(err)
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
