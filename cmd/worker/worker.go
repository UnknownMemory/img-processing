package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Worker struct {
	RMQ    string
	logger *log.Logger
}

func NewWorker(RMQ string, logger *log.Logger) *Worker {
	return &Worker{
		RMQ:    RMQ,
		logger: logger,
	}
}

func (worker *Worker) Connect() {
	conn, err := amqp.Dial(worker.RMQ)
	failOnError(err, "Failed to connect to RabbitMQ")

	defer func(conn *amqp.Connection) {
		err := conn.Close()
		failOnError(err, "Failed to close AMQP connection")
	}(conn)

	ch, _ := conn.Channel()
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		failOnError(err, "Failed to close connection channel")
	}(ch)

	q, err := ch.QueueDeclare("image", true, false, false, false, nil)
	if err != nil {
		worker.logger.Panicf("Failed to declare queue")
	}

	messages, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		worker.logger.Panicf("Failed to register a consumer")
	}

	forever := make(chan bool)

	go worker.Receiver(messages)
	worker.logger.Println("Waiting for messages")
	<-forever
}

func (worker *Worker) Receiver(messages <-chan amqp.Delivery) {
	for message := range messages {
		println(message.Body)
	}
}
