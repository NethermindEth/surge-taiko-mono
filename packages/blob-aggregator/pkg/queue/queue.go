package queue

import (
	"context"
	"errors"
	"sync"

	"github.com/taikoxyz/taiko-mono/packages/blob-aggregator/pkg/types"
)

var (
	ErrClosed = errors.New("queue connection closed")
)

type Queue interface {
	Close()
	Ack(ctx context.Context, msg Message) error
	Publish(ctx context.Context, proposal types.QueueProposalRequestBody) error
	Subscribe(ctx context.Context, msgChan chan<- Message, wg *sync.WaitGroup) error
}

type Message struct {
	Proposal types.QueueProposalRequestBody
	Internal interface{}
}

type NewQueueOpts struct {
	Username      string
	Password      string
	Host          string
	Port          string
	PrefetchCount uint64
}
