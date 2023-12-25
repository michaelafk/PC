package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	const MinimBoti = 20 // Xifra a partir de la qual el tesorer acabarà

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/") //Conecta al servidor
	failOnError(err, "Tesorer: Error establint la connexió amb RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel() //Obre un canal
	failOnError(err, "Tesorer: Error al obrir un canal")
	defer ch.Close()

	deposits, err := ch.QueueDeclare( //Cua deposits
		"deposits", // name
		false,      // durable (Si el servidor cau la cua perdura amb combinació amb els missatges persistents)
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Tesorer: Error al declarar la cua deposits")

	depMsgs, err := ch.Consume( //Missatges de la cua deposits
		deposits.Name, // queue
		"",            // consumer
		false,         // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	failOnError(err, "Tesorer: Error al declarar els missatges de la cua deposits")

	balances, err := ch.QueueDeclare( //Cua banlances
		"balances", // name
		false,      // durable (Si el servidor cau la cua perdura amb combinació amb els missatges persistents)
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Tesorer: Error al declarar la cua balances")

	err = ch.ExchangeDeclare( // Intercanvi de tipus fanout per fer acabar als clients
		"robatori", // Nom de l'intercanvi
		"fanout",   // Tipus d'intercanvi
		true,       // Durable
		false,      // No esborrable
		false,      // Exclusiu
		false,      // No esperar a l'intercanvi ja declarat
		nil,        // Arguments addicionals
	)
	failOnError(err, "Tesorer: Error al declarar l'intercamvi robatori")

	fmt.Println("El tresorer és al despatx. El botí mínim és: " + strconv.Itoa(MinimBoti))
	log.Printf(" [*] Esperant clients")
	total := 0 // Balanç total quan arriba a un MinimBoti el Tesorer roba el que hi ha
	robat := false

	for d := range depMsgs {
		fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
		body := string(d.Body)
		parts := strings.Split(body, ";") // ; Separador nom i valor diposit
		if len(parts) != 2 {
			log.Printf("Tresorer: Error en el missatge rebut: %s", body)
			continue
		}
		nom := parts[0]
		n := parts[1]
		fmt.Println("Operació rebuda: " + n + " del client: " + nom)
		d.Ack(true)
		valor, _ := strconv.Atoi(n)
		if total+valor >= 0 {
			total += valor
		} else {
			fmt.Println("OPERACIÓ NO PERMESA NO HI HA FONS")
		}
		fmt.Println("Balanç: " + strconv.Itoa(total))
		if total >= MinimBoti {
			total = 0 // Roba el diposit
			robat = true
			fmt.Println("El Tresorer decideix robar el diposit i tancar el despatx")
		}
		err = ch.Publish( // Pulicar el missatge a la cua dels balanços
			"",            // Nom de l'intercamvi (buit ja que publiquem directament a la cua)
			balances.Name, // Nom de la cua
			false,         // No esperar confirmació
			false,         // No eliminar missatges no lliurats
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(strconv.Itoa(total)),
			})
		failOnError(err, "Tresorer: Error en publicar a balances")
		fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
		if robat {
			fmt.Println("El Tresorer se'n va\n\n")
			time.Sleep(time.Duration(rand.Intn(9000)) * time.Millisecond)
			err = ch.Publish( // Publicar el missatge a l'intercanvi robatori de tipus fanout
				"robatori", // Nom de l'intercanvi
				"",         // La clau de l'intercanvi (fanout no utilitza claus)
				false,      // Esperar confirmació
				false,      // No eliminar missatges no entregats
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte("L'oficina acaba de tancar"),
				},
			)
			failOnError(err, "Tresorer: Error en publicar a robatori")
			fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
			// Esborrar totes les cues
			queueNames := []string{"deposits", "balances"}
			for _, queueName := range queueNames {
				_, err := ch.QueueDelete(
					queueName, // Nom de la cua
					false,     // No esborrar només si no té consumidors
					false,     // No esborrar tot i amb consumidors actius
					false,     // No esperar a la cua ja eliminada
				)
				failOnError(err, fmt.Sprintf("Failed to delete queue '%s'", queueName))
				fmt.Printf("Cua '%s' esborrada amb èxit.\n", queueName)
			}
			break
		}
	}
}
