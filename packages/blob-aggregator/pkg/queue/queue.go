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
	Publish(ctx context.Context, proposal types.QueueProposalRequestBody) error
	Subscribe(ctx context.Context, msgChan chan<- types.QueueProposalRequestBody, wg *sync.WaitGroup) error
}

type NewQueueOpts struct {
	Username      string
	Password      string
	Host          string
	Port          string
	PrefetchCount uint64
}
