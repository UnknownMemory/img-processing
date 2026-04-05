package rabbitmq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/unknownmemory/img-processing/internal/aws"
	db "github.com/unknownmemory/img-processing/internal/database"
	process "github.com/unknownmemory/img-processing/internal/proc"
	"github.com/unknownmemory/img-processing/internal/shared"
)

type RabbitMQ struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	logger *log.Logger
	db     *pgxpool.Pool
}

func NewWorker(RMQ string, logger *log.Logger, db *pgxpool.Pool) (*RabbitMQ, error) {
	conn, err := amqp.Dial(RMQ)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &RabbitMQ{
		conn:   conn,
		ch:     ch,
		logger: logger,
		db:     db,
	}, nil
}

func (worker *RabbitMQ) Close() {
	if worker.ch != nil {
		if err := worker.ch.Close(); err != nil {
			worker.logger.Printf("failed to close channel: %s", err)
		}
	}

	if worker.conn != nil {
		if err := worker.conn.Close(); err != nil {
			worker.logger.Printf("failed to close connection: %s", err)
		}
	}
}

func (worker *RabbitMQ) Listen() {

	q, err := worker.ch.QueueDeclare("image", true, false, false, false, nil)
	if err != nil {
		worker.logger.Panicf("Failed to declare queue")
	}

	err = worker.ch.Qos(1, 0, false)
	failOnError(err, "Failed to set QoS")

	messages, err := worker.ch.Consume(q.Name, "", false, false, false, false, nil)
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
		transform, mime, err := process.Transform(object, data.Transformations)
		if err != nil {
			return
		}

		transformKey := fmt.Sprintf("%v/%s/image", message.Headers["userId"], message.Headers["uuid"])
		_, err = awsCli.Upload(transformKey, bytes.NewReader(transform), mime)
		if err != nil {
			return
		}

		imageUUID, err := uuid.Parse(message.Headers["uuid"].(string))
		if err != nil {
			return
		}

		userId := message.Headers["userId"].(string)
		uId, err := strconv.ParseInt(userId, 10, 64)
		if err != nil {
			return
		}

		q := db.New(worker.db)
		transformQuery := &db.UpdateTransformParams{
			Status: "completed",
			Mime:   pgtype.Text{String: mime, Valid: true},
			Uuid:   pgtype.UUID{Bytes: imageUUID, Valid: true},
			UserID: pgtype.Int8{Int64: uId, Valid: true},
		}

		err = q.UpdateTransform(context.Background(), *transformQuery)
		if err != nil {
			return
		}

		err = message.Ack(false)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func (worker *RabbitMQ) Send(queueName string, data interface{}, userId string, transformUUID string) {
	body, err := json.Marshal(data)
	if err != nil {
		failOnError(err, "Failed to marshal data")
	}

	q, err := worker.ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		worker.logger.Panicf("Failed to declare queue")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = worker.ch.PublishWithContext(ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			Headers:      amqp.Table{"userId": userId, "uuid": transformUUID},
			DeliveryMode: amqp.Persistent,
		})
	failOnError(err, "Failed to publish a message")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", err, msg)
	}
}
