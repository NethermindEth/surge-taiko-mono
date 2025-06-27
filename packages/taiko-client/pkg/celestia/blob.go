package celestia

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/celestiaorg/go-square/merkle"
	"github.com/celestiaorg/go-square/v2/inclusion"
	"github.com/celestiaorg/go-square/v2/share"
)

const (
	// https://github.com/celestiaorg/celestia-app/blob/main/pkg/appconsts/v4/app_consts.go#L17
	CelestiaBlobNamespace         = "blob"
	SubtreeRootThreshold  int     = 64
	DefaultGasPrice       float64 = -1.0
	DefaultMinGasPrice            = 0.002
	DefaultMaxGasPrice            = DefaultMinGasPrice * 100
)

type Commitment []byte

type Blob struct {
	*share.Blob `json:"blob"`

	Commitment Commitment `json:"commitment"`

	index int
}

type SubmitOptions struct {
	signerAddress     string
	keyName           string
	gasPrice          float64
	isGasPriceSet     bool
	maxGasPrice       float64
	gas               uint64
	priority          int
	feeGranterAddress string
}

type CelestiaBlobHandler struct {
	Submit func(context.Context, []*Blob, *SubmitOptions) (uint64, error)
}

func NewBlobV0(namespace share.Namespace, data []byte) (*Blob, error) {
	return NewBlob(share.ShareVersionZero, namespace, data, nil)
}

func NewBlob(shareVersion uint8, namespace share.Namespace, data, signer []byte) (*Blob, error) {
	if err := namespace.ValidateForBlob(); err != nil {
		return nil, fmt.Errorf("invalid user namespace: %w", err)
	}

	shareBlob, err := share.NewBlob(namespace, data, shareVersion, signer)
	if err != nil {
		return nil, err
	}

	commitment, err := inclusion.CreateCommitment(shareBlob, merkle.HashFromByteSlices, SubtreeRootThreshold)
	if err != nil {
		return nil, err
	}

	return &Blob{Blob: shareBlob, Commitment: commitment, index: -1}, nil
}

type jsonBlob struct {
	Namespace    []byte     `json:"namespace"`
	Data         []byte     `json:"data"`
	ShareVersion uint8      `json:"share_version"`
	Commitment   Commitment `json:"commitment"`
	Signer       []byte     `json:"signer,omitempty"`
	Index        int        `json:"index"`
}

func (b *Blob) MarshalJSON() ([]byte, error) {
	blob := &jsonBlob{
		Namespace:    b.Namespace().Bytes(),
		Data:         b.Data(),
		ShareVersion: b.ShareVersion(),
		Commitment:   b.Commitment,
		Signer:       b.Signer(),
		Index:        b.index,
	}
	return json.Marshal(blob)
}

func (b *Blob) UnmarshalJSON(data []byte) error {
	var jsonBlob jsonBlob
	err := json.Unmarshal(data, &jsonBlob)
	if err != nil {
		return err
	}

	ns, err := share.NewNamespaceFromBytes(jsonBlob.Namespace)
	if err != nil {
		return err
	}

	blob, err := NewBlob(jsonBlob.ShareVersion, ns, jsonBlob.Data, jsonBlob.Signer)
	if err != nil {
		return err
	}

	blob.Commitment = jsonBlob.Commitment
	blob.index = jsonBlob.Index
	*b = *blob
	return nil
}

func NewSubmitOptions() *SubmitOptions {
	return &SubmitOptions{
		gasPrice:    DefaultGasPrice,
		maxGasPrice: DefaultMaxGasPrice,
	}
}
