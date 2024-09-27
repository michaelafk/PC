// Noé Fernández Simmank
package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Empty struct{}

func client(ch *amqp.Channel, ch2 *amqp.Channel, colaDiposits amqp.Queue, colaBalanç amqp.Queue, espera amqp.Queue, done chan Empty, nombre string, numOperaciones int) {

	fmt.Println(fmt.Sprintf("Hola el meu nom és: %s", nombre))
	fmt.Println(fmt.Sprintf("%s vol fer %d operacions", nombre, numOperaciones))

	// abre el consumo para el mensaje
	msgT, err := ch2.Consume(
		espera.Name, // name
		"Cliente",   // id
		false,
		false,
		false,
		false,
		nil,
	)
	checkFail(err, "Error al abrir cola de espera")

	for i := 0; i < numOperaciones; i++ {

		//se obtiene el valor de la opercion
		var operacion int
		operacion = RandInt(-10, 10)

		fmt.Println(fmt.Sprintf("%s operacion %d: %d", nombre, i+1, operacion))

		message := fmt.Sprintf("%d, %s", operacion, nombre)
		//se publica la operacion
		err2 := ch.Publish(
			"",
			colaDiposits.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(message),
			})
		checkFail(err2, "Error al publicar un mensaje")

		fmt.Println("Operació sol·licitada...")
		//espera para simular que la contestación no es inmediata
		time.Sleep(time.Second)
		timeout := time.After(time.Millisecond)

		//entonces comprueva si el tesorero se ha marchado, si es asi acaba, si no continua
		select {
		case msg := <-msgT:
			fmt.Println(fmt.Sprintf("El tresorer a dit: %s", string(msg.Body)))
			fmt.Println("Jo també me'n vaig idò!")
			done <- Empty{}
			os.Exit(0)
		case <-timeout:
			//recibe el mensaje de balance
			msg, r, err := ch.Get("colaBalanç", true)
			if !r {
				checkFail(err, "Error al coger una pieza del plato")
			}
			balance, err := strconv.Atoi(string(msg.Body))
			if balance == 0 && operacion > 0 {
				fmt.Println("EL TRESORER SEMBLA SOSPITÓS")
			} else if balance == 0 && operacion < 0 {
				fmt.Println("NO HI HA SALDO")
			} else if balance > 0 && operacion > 0 {
				fmt.Println("OPERACIÓ CORRECTE")
			} else if balance > 0 && operacion < 0 {
				fmt.Println("ES FARÀ EL REINTEGRE SI HI HA SALDO")
			}
			fmt.Println(fmt.Sprintf("Balanç actual: %d", balance))
			fmt.Println(fmt.Sprintf("%d----------------------------------------", i+1))
		}
	}

	// envia el mensaje por el canal para hacer saber que ha terminado
	done <- Empty{}
}

// genera un número random
func RandInt(lower, upper int) int {
	rand.Seed(time.Now().UnixNano())
	rng := upper - lower
	ret := rand.Intn(rng) + lower
	// en el caso de que el valor sea 0 se repite hasta que no lo sea
	for {
		rand.Seed(time.Now().UnixNano())
		ret = rand.Intn(rng) + lower
		if ret != 0 {
			break
		}
	}
	return ret
}

func main() {

	args := os.Args

	var nombre string
	nombre = args[1]

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
	colaBalanç, err := ch.QueueDeclare(
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

	// cola de espera
	espera, err := ch2.QueueDeclare(
		"",
		false,
		false,
		true, // exclusive
		false,
		nil,
	)
	checkFail(err, "Error al crear una cola ")

	// vincular cola con el exchange
	err = ch2.QueueBind(
		espera.Name, // queue name
		"",
		"espera", // exchange
		false,
		nil,
	)
	checkFail(err, "Error bind-cola")

	// randomizar la semilla del número aleatorio
	rand.Seed(time.Now().UnixNano())
	//calculo de las operaciones que realizara el cliente
	var numOperaciones int
	numOperaciones = rand.Intn(5) + 1
	//lanza el proceso cliente
	go client(ch, ch2, colaDiposits, colaBalanç, espera, done, nombre, numOperaciones)

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
