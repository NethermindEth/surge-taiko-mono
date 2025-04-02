package aggregator

import (
	"log/slog"

	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/taikoxyz/taiko-mono/packages/blob-aggregator/pkg/queue"
	taikoEncoding "github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/encoding"
)

func (agg *Aggregator) buildBlobs(msgBatch []*queue.Message) ([]*eth.Blob, error) {
	var (
		i         uint
		blobBytes hexutil.Bytes = make(hexutil.Bytes, 0, eth.BlobSize)
		start     uint
	)

	var blobs []*eth.Blob

	for {
		byteSize := min(uint(len(msgBatch[i].Proposal.TxList[start:])), eth.BlobSize)
		blobBytes = append(blobBytes, msgBatch[i].Proposal.TxList[start:start+byteSize]...)
		start += byteSize

		// If the current tx list is completely consumed
		if start == uint(len(msgBatch[i].Proposal.TxList)) {
			start = 0
			i++
		}

		// If enough data is collected for the blob, OR we have consumed all tx lists
		if len(blobBytes) == eth.BlobSize || i == uint(len(msgBatch)) {
			var blob = &eth.Blob{}
			err := blob.FromData(blobBytes)
			if err != nil {
				return nil, err
			}
			blobs = append(blobs, blob)

			blobBytes = hexutil.Bytes{}
		}

		// If all the transactions have been consumed
		if i == uint(len(msgBatch)) {
			break
		}
	}

	return blobs, nil
}

func (agg *Aggregator) buildBlobParamsForBatchedProposal() []*taikoEncoding.BlobParams {
	var (
		byteOffset     uint
		blobNumber     uint
		blobParamsList []*taikoEncoding.BlobParams
	)

	for _, msg := range agg.batch {
		txListLen := len(msg.Proposal.TxList)
		blobsUsed := (byteOffset + uint(txListLen)) / uint(eth.BlobSize)

		blobParams := taikoEncoding.BlobParams{
			FirstBlobIndex: uint8(blobNumber),
			NumBlobs:       uint8(txListLen / eth.BlobSize),
			ByteOffset:     uint32(byteOffset),
			ByteSize:       uint32(txListLen),
		}

		byteOffset = (byteOffset + uint(txListLen)) % uint(eth.BlobSize)
		blobNumber += blobsUsed

		// Store the blob params locally and remove the now redundant proposal message from rabbit mq
		blobParamsList = append(blobParamsList, &blobParams)
		agg.queue.Ack(agg.ctx, *msg)

		slog.Info("aggregated new proposal", "proposal", msg.Proposal)
	}

	return blobParamsList
}
