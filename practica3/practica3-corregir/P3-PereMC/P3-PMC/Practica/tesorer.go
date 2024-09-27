// @Author: Pere Martínez Cifre
package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

const (
	MinimBoti = 20 // La quantitat mínima a conseguir per robar-ho
)

var balanç = 0 // Variable global per controlar el balanç del compte

// Imprimir una línia en format log a consola
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// Connectar amb RabbitMQ
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Error creant la conexió amb RabbitMQ")
	defer conn.Close()

	// Obrir el canal dins rabbitMQ
	ch, err := conn.Channel()
	failOnError(err, "Error obrint el canal")
	defer ch.Close()

	// Crear o connectar (Si ja s'ha creat) l'exchange necessari per
	// el missatge de despedida de tipus Fanout. Cada client crearà
	// una cua única que també vincularà a aquest exchange. D'aquesta
	// manera, la quantitat de cues pot ser dinàmica.
	err = ch.ExchangeDeclare(
		"fanoutExchange", // nom de l'exchange
		"fanout",         // tipus d'exchange
		false,            // durable?
		true,             // auto-borrar?
		false,            // interna?
		false,            // sense-espera?
		nil,              // arguments
	)
	failOnError(err, "Error creant l'exchange fanout")

	// Crear la cua Balances, aquesta cua l'utilitzaràn els clients per indicar
	// les diferentes operacions que desitgen fer. El tesorer les cullirà d'una
	// en una per processar-les.
	qBalances, err := ch.QueueDeclare(
		"Balances", // nom de la cua
		false,      // durable?
		false,      // auto-borrar?
		false,      // exclusiva?
		false,      // sense-espera?
		nil,        // arguments
	)
	failOnError(err, "Error creant la cua Balances")

	// Crear la cua Diposits, aquesta cua l'utilitzarà el tesorer per publicar
	// el balanç del compte després de cada operació processada de Balances.
	// Els clientes culliràn els missatges i els compararàn amb l'operació
	// publicada per ells a la cua Balances.
	qDiposits, err := ch.QueueDeclare(
		"Diposits", // nom de la cua
		false,      // durable?
		false,      // auto-borrar?
		false,      // exclusiva?
		false,      // sense-espera?
		nil,        // arguments
	)
	failOnError(err, "Error creant la cua Diposits")

	// Crear relació de tipus consumidor amb la cua Balances. Això permet al
	// tesorer anar collint les diferents operacions que volen fer els clients.
	missatgesBalanç, err := ch.Consume(
		qBalances.Name, // nom de la cua
		"",             // identificador del consumidor
		true,           // auto-ack?
		false,          // exclusiva?
		false,          // no-local?
		false,          // sense-espera?
		nil,            // args
	)
	failOnError(err, "Error registrant el consumidor per Balances")

	forever := make(chan bool) //Canal utilitzat per evitar que el programa acabi

	// --------------- Funcionament --------------- //
	go func() {
		// Per cada operació sol·licitada per algun client
		for missatge := range missatgesBalanç {
			// El client envia l'operacio en format "valorOperacio|nomClient"
			// Amb aquest split tenim
			// [0] valorOperacio o "" (Si la cua ha tancat)
			// [1] nomClient o OutOfRange (Si la cua ha tancat)
			dadesRebudes := strings.Split(string(missatge.Body), "|")
			fmt.Printf(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>\n")

			// En cas que la cua hagi tancat, ja no s'ha de processar cap altre missatge
			if len(dadesRebudes) < 2 {
				break
			}

			// Operació rebuda
			fmt.Printf("Operació rebuda: " + dadesRebudes[0] + " del client: " + dadesRebudes[1] + "\n")

			// Collim el valorOperacio i el convertim a enter
			valorOperacio, err := strconv.Atoi(dadesRebudes[0])
			failOnError(err, "Error durant la conversió a enter")

			// Processam l'operació, si és vàlida, l'aplicam al balanç del compte
			balançTemporal := balanç + valorOperacio
			if balançTemporal < 0 {
				fmt.Printf("OPERACIÓ NO PERMESA NO HI HA FONS\n")
			} else {
				balanç = balançTemporal
			}

			fmt.Printf("Balanç: %d\n", balanç) // Indicam el balanç a consola

			// Preparam per publicar a la cua Diposits el nou balanç
			body := strconv.Itoa(balanç)
			ctx := context.TODO()
			err = ch.PublishWithContext(
				ctx,            // contexte (genèric en el nostre cas)
				"",             // exchange (predeterminat per arribar a una cua)
				qDiposits.Name, // cua objectiu
				false,          // isMandatory?
				false,          // isImmediate?
				amqp091.Publishing{
					ContentType: "text/plain",
					Body:        []byte(body),
				})
			failOnError(err, "Error publicant a la cua Diposits")

			// Si el balanç ja és vàlid per robar
			if balanç >= MinimBoti {
				fmt.Printf("El tresorer decideix robar el dipòsit i tancar el despatx\n")

				// Enviam a totes les cues d'informació (les qComunicacio) de
				// cada client el missatge de que es tanca l'oficina.
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				body := "L'oficina acaba de tancar"
				err = ch.PublishWithContext(ctx,
					"fanoutExchange", // exchange
					"",               // routing key
					false,            // mandatory
					false,            // immediate
					amqp091.Publishing{
						ContentType: "text/plain",
						Body:        []byte(body),
					})
				failOnError(err, "Error publicant l'avís en Fanout")

				forever <- false //Eliminam el bloqueig del fil per permetre acabar
			}

			fmt.Printf("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<\n")
		}
	}()

	//Indicam l'inici del tresorer
	fmt.Printf("El tresorer és al despatx. El botí mínim és: " + strconv.Itoa(MinimBoti) + "\n")

	<-forever

	//Indicam que el tresorer marxa
	fmt.Printf("El Tresorer se'n va\n\n")

	//Iniciam l'eliminació de les cues Balances i Diposits
	delete, err := ch.QueueDelete("Balances", false, false, false)
	failOnError(err, "Error esborrant la cua Balances")
	log.Printf("Cua 'Balances' esborrada amb èxit.\n")

	//Per evitar l'error de que no s'utilitza el primer valot de delete.
	//No fa res real
	delete = delete + 1

	delete, err = ch.QueueDelete("Diposits", false, false, false)
	failOnError(err, "Error esborrant la cua Dipòsits")
	log.Printf("Cua 'Depòsits' esborrada amb èxit.\n")
}
