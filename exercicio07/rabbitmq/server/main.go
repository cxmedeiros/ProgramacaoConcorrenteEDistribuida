package main

import (
	"Vituriano/sorters"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func QuickSort(arr []int64) []int64 {
	intArray := arr
	sorters.QuickSortAsync(intArray, 0, int64(len(intArray)-1))
	sortedArray := intArray
	return sortedArray
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,
		0,
		false,
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {

			var arr []int64
			err := json.Unmarshal(d.Body, &arr)
			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
				d.Ack(false)
				continue
			}

			sortedArr := QuickSort(arr)

			sortedArrJSON, err := json.Marshal(sortedArr)
			if err != nil {
				log.Printf("Error encoding JSON: %s", err)
				d.Ack(false)
				continue
			}

			err = ch.Publish(
				"",        // exchange
				d.ReplyTo, // routing key (ReplyTo queue from the request)
				false,     // mandatory
				false,
				amqp.Publishing{
					ContentType:   "application/json",
					Body:          sortedArrJSON,
					CorrelationId: d.CorrelationId, // same correlation ID
				})
			if err != nil {
				log.Printf("Error publishing response: %s", err)
			}

			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
