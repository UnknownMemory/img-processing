package rabbitmq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	amqp "github.com/rabbitmq/amqp091-go"
	db "github.com/unknownmemory/img-processing/internal/database"
	"github.com/unknownmemory/img-processing/internal/shared"
)

type Storage interface {
	GetObject(key string) ([]byte, error)
	Upload(key string, body io.ReadSeeker, mime string) error
}

type ImageProcessor interface {
	Transform(object []byte, operations shared.Transformations) ([]byte, string, error)
}

type TransformRepository interface {
	UpdateTransform(ctx context.Context, arg db.UpdateTransformParams) error
}

type RabbitMQ struct {
	conn    *amqp.Connection
	ch      *amqp.Channel
	logger  *log.Logger
	storage Storage
	proc    ImageProcessor
	db      TransformRepository
}

func NewWorker(RMQ string, logger *log.Logger, storage Storage, proc ImageProcessor, repo TransformRepository) (*RabbitMQ, error) {
	conn, err := amqp.Dial(RMQ)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &RabbitMQ{
		conn:    conn,
		ch:      ch,
		logger:  logger,
		storage: storage,
		proc:    proc,
		db:      repo,
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

func (worker *RabbitMQ) Listen(workerPool int) {

	q, err := worker.ch.QueueDeclare("image", true, false, false, false, nil)
	if err != nil {
		worker.logger.Panicf("Failed to declare queue")
	}

	err = worker.ch.Qos(2, 0, false)
	failOnError(err, "Failed to set QoS")

	messages, err := worker.ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		worker.logger.Panicf("Failed to register a consumer")
	}

	for i := 0; i < workerPool; i++ {
		go worker.Receiver(i+1, messages)
	}
	worker.logger.Printf("Waiting for messages with %d workers", workerPool)

	forever := make(chan bool)
	<-forever
}

func (worker *RabbitMQ) Receiver(workerID int, messages <-chan amqp.Delivery) {
	for message := range messages {
		worker.logger.Printf("Worker %d processing message", workerID)

		err := worker.handleMessage(message)
		if err != nil {
			worker.logger.Printf("Failed to process message: %d", err)
			nackErr := message.Nack(false, false)
			if nackErr != nil {
				worker.logger.Printf("Failed to nack message: %d", nackErr)
			}
			continue
		}

		err = message.Ack(false)
		if err != nil {
			worker.logger.Printf("Failed to acknowledge message: %d", err)
			continue
		}

		worker.logger.Printf("Worker %d completed message", workerID)
	}
}

func (worker *RabbitMQ) handleMessage(message amqp.Delivery) error {
	data := &shared.ImageTransform{}
	_ = json.Unmarshal(message.Body, &data)

	key := fmt.Sprintf("%v/%s/original", message.Headers["userId"], data.ImageID)
	object, err := worker.storage.GetObject(key)
	if err != nil {
		return fmt.Errorf("failed to get object from S3: %d", err)
	}

	transform, mime, err := worker.proc.Transform(object, data.Transformations)
	if err != nil {
		return fmt.Errorf("failed to transform image: %d", err)
	}

	transformKey := fmt.Sprintf("%v/%s/image", message.Headers["userId"], message.Headers["uuid"])
	err = worker.storage.Upload(transformKey, bytes.NewReader(transform), mime)
	if err != nil {
		return fmt.Errorf("failed to upload transformed image to S3: %d", err)
	}

	imageUUID, err := uuid.Parse(message.Headers["uuid"].(string))
	if err != nil {
		return fmt.Errorf("failed to parse image UUID: %d", err)
	}

	userId := message.Headers["userId"].(string)
	uId, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse user ID: %d", err)
	}

	transformQuery := db.UpdateTransformParams{
		Status: "completed",
		Mime:   pgtype.Text{String: mime, Valid: true},
		Uuid:   pgtype.UUID{Bytes: imageUUID, Valid: true},
		UserID: pgtype.Int8{Int64: uId, Valid: true},
	}

	err = worker.db.UpdateTransform(context.Background(), transformQuery)
	if err != nil {
		return fmt.Errorf("failed to update transform status: %d", err)
	}

	return nil
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
