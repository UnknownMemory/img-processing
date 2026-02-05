package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/unknownmemory/img-processing/internal/aws"
	process "github.com/unknownmemory/img-processing/internal/proc"
	"github.com/unknownmemory/img-processing/internal/shared"
)

type RabbitMQ struct {
	RMQ    string
	logger *log.Logger
}

func NewWorker(RMQ string, logger *log.Logger) *RabbitMQ {
	return &RabbitMQ{
		RMQ:    RMQ,
		logger: logger,
	}
}

func (worker *RabbitMQ) Listen() {
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

func (worker *RabbitMQ) Receiver(messages <-chan amqp.Delivery) {
	for message := range messages {
		data := &shared.ImageTransform{}
		_ = json.Unmarshal(message.Body, &data)

		key := fmt.Sprintf("%v/%s/original", message.Headers["userId"], data.ImageID)
		awsCli := aws.NewS3Client()
		object, err := awsCli.GetObject(key)
		if err != nil {
			return
		}
		transform, err := process.Transform(object, data.Transformations)
		if err != nil {
			return
		}
	}
}

func (worker *RabbitMQ) Send(queueName string, data interface{}, userId string) {
	body, err := json.Marshal(data)
	if err != nil {
		failOnError(err, "Failed to marshal data")
	}

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

	q, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		worker.logger.Panicf("Failed to declare queue")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Headers:     amqp.Table{"userId": userId},
		})
	failOnError(err, "Failed to publish a message")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", err, msg)
	}
}
