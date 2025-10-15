package celestia

import (
	"context"
)

const (
	CelestiaStateNamespace = "state"
)

type Balance struct {
	Amount string `json:"amount"`

	Denomination string `json:"denom"`
}

type CelestiaStateHandler struct {
	Balance func(ctx context.Context) (Balance, error)
}
