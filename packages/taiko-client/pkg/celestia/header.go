package celestia

import (
	"context"
)

const (
	CelestiaHeaderNamespace = "header"
)

type CelestiaHeaderHandler struct {
	NetworkHead func(ctx context.Context) error
}
