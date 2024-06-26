package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{message}", messageHandler)
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func messageHandler(w http.ResponseWriter, r *http.Request) {
	message := r.PathValue("message")

	conn, err := amqp.Dial("amqp://guest:guest@34.94.104.70:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	failOnError(err, "Failed to publish a message")
	fmt.Fprintf(w, "'%s' message sent", message)
}
