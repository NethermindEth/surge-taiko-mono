package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/taikoxyz/taiko-mono/packages/blob-aggregator/pkg/queue"
	"github.com/taikoxyz/taiko-mono/packages/blob-aggregator/pkg/types"
)

type RabbitMQ struct {
	conn      *amqp.Connection
	ch        *amqp.Channel
	queueName string

	connErrCh chan *amqp.Error
	chErrCh   chan *amqp.Error

	subscriptionCtx    context.Context
	subscriptionCancel context.CancelFunc

	opts queue.NewQueueOpts
}

func NewRabbitMQ(opts queue.NewQueueOpts) (*RabbitMQ, error) {
	slog.Info("dialing rabbitmq connection")

	r := &RabbitMQ{
		opts: opts,
	}

	err := r.connect()
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *RabbitMQ) connect() error {
	slog.Info("connecting to rabbitmq")

	conn, err := amqp.DialConfig(
		fmt.Sprintf(
			"amqp://%v:%v@%v:%v/",
			r.opts.Username,
			r.opts.Password,
			r.opts.Host,
			r.opts.Port,
		),
		amqp.Config{
			Heartbeat: 1 * time.Second,
		})
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	if err := ch.Qos(int(r.opts.PrefetchCount), 0, false); err != nil {
		return err
	}

	r.conn = conn
	r.ch = ch

	slog.Info("connected to rabbitmq")

	return nil
}

func (r *RabbitMQ) Publish(ctx context.Context, proposal types.QueueProposalRequestBody) error {
	slog.Info("Publishing message to RabbitMQ", "queue", r.queueName)

	body, err := json.Marshal(proposal)
	if err != nil {
		return err
	}

	err = r.ch.PublishWithContext(ctx,
		"",
		r.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}

	slog.Info("Message published successfully")
	return nil
}

func (r *RabbitMQ) Subscribe(ctx context.Context, msgChan chan<- types.QueueProposalRequestBody, wg *sync.WaitGroup) error {
	wg.Add(1)
	defer wg.Done()

	slog.Info("Starting message consumer", "queue", r.queueName)

	msgs, err := r.ch.Consume(
		r.queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for {
		select {
		case <-r.subscriptionCtx.Done():
			defer r.Close()
			slog.Info("Subscription context cancelled")
			return nil
		case <-ctx.Done():
			defer r.Close()
			slog.Info("Consumer context cancelled")
			return nil
		case err := <-r.connErrCh:
			slog.Error("RabbitMQ connection closed", "error", err)
			return err
		case err := <-r.chErrCh:
			slog.Error("RabbitMQ channel closed", "error", err)
			return err
		case d, ok := <-msgs:
			if !ok {
				slog.Error("Message channel closed")
				return nil
			}

			var proposal types.QueueProposalRequestBody
			err := json.Unmarshal(d.Body, &proposal)
			if err != nil {
				slog.Error("Failed to parse message", "error", err)
				_ = d.Nack(false, false)
				continue
			}

			slog.Info("Received message", "txDest", proposal.TxDest)
			msgChan <- proposal

			_ = d.Ack(false)
		}
	}
}

func (r *RabbitMQ) Close() {
	if r.subscriptionCancel != nil {
		r.subscriptionCancel()
	}

	if err := r.ch.Close(); err != nil && err != amqp.ErrClosed {
		slog.Error("Error closing RabbitMQ channel", "error", err)
	}

	if err := r.conn.Close(); err != nil && err != amqp.ErrClosed {
		slog.Error("Error closing RabbitMQ connection", "error", err)
	}

	slog.Info("RabbitMQ connection closed")
}
