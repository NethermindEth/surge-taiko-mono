package transaction

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum-optimism/optimism/op-service/txmgr"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"

	"github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/encoding"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/metadata"
	pacayaBindings "github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/pacaya"
	surgeBindings "github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings/surge"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
	proofProducer "github.com/taikoxyz/taiko-mono/packages/taiko-client/prover/proof_producer"
)

var (
	ErrUnretryableSubmission = errors.New("unretryable submission error")
)

// TxBuilder will build a transaction with the given nonce.
type TxBuilder func(txOpts *bind.TransactOpts) (*txmgr.TxCandidate, error)

// ProveBatchesTxBuilder is responsible for building ProveBatches transactions.
type ProveBatchesTxBuilder struct {
	rpc                *rpc.Client
	pacayaInboxAddress common.Address
	shastaInboxAddress common.Address
	proverSetAddress   common.Address
}

// NewProveBatchesTxBuilder creates a new ProveBatchesTxBuilder instance.
func NewProveBatchesTxBuilder(
	rpc *rpc.Client,
	pacayaInboxAddress common.Address,
	shastaInboxAddress common.Address,
	proverSetAddress common.Address,
) *ProveBatchesTxBuilder {
	return &ProveBatchesTxBuilder{rpc, pacayaInboxAddress, shastaInboxAddress, proverSetAddress}
}

// BuildProveBatchesPacaya creates a new TaikoInbox.ProveBatches transaction.
func (a *ProveBatchesTxBuilder) BuildProveBatchesPacaya(batchProof *proofProducer.BatchProofs) TxBuilder {
	return func(txOpts *bind.TransactOpts) (*txmgr.TxCandidate, error) {
		var (
			data        []byte
			to          common.Address
			err         error
			metas       = make([]metadata.TaikoProposalMetaData, len(batchProof.ProofResponses))
			transitions = make([]pacayaBindings.ITaikoInboxTransition, len(batchProof.ProofResponses))
			subProofs   []encoding.SubProofPacaya
			batchIDs    = make([]uint64, len(batchProof.ProofResponses))
		)
		for i, proof := range batchProof.ProofResponses {
			metas[i] = proof.Meta
			transitions[i] = pacayaBindings.ITaikoInboxTransition{
				ParentHash: proof.Opts.PacayaOptions().Headers[0].ParentHash,
				BlockHash:  proof.Opts.PacayaOptions().Headers[len(proof.Opts.PacayaOptions().Headers)-1].Hash(),
				StateRoot:  proof.Opts.PacayaOptions().Headers[len(proof.Opts.PacayaOptions().Headers)-1].Root,
			}
			batchIDs[i] = proof.Meta.Pacaya().GetBatchID().Uint64()
			log.Info(
				"Build batch proof submission transaction",
				"batchID", batchIDs[i],
				"parentHash", common.Bytes2Hex(transitions[i].ParentHash[:]),
				"blockHash", common.Bytes2Hex(transitions[i].BlockHash[:]),
				"stateRoot", common.Bytes2Hex(transitions[i].StateRoot[:]),
				"startBlockID", proof.Opts.PacayaOptions().Headers[0].Number,
				"endBlockID", proof.Opts.PacayaOptions().Headers[len(proof.Opts.PacayaOptions().Headers)-1].Number,
				"gasLimit", txOpts.GasLimit,
			)
		}
		log.Info(
			"Verifier information",
			"ProofType1", batchProof.ProofType1,
			"Verifier1", batchProof.Verifier1,
			"Proof1", common.Bytes2Hex(batchProof.BatchProof1),
			"ProofType2", batchProof.ProofType2,
			"Verifier2", batchProof.Verifier2,
			"Proof2", common.Bytes2Hex(batchProof.BatchProof2),
		)

		subProofs = []encoding.SubProofPacaya{
			{
				Verifier: batchProof.Verifier1,
				Proof:    batchProof.BatchProof1,
			},
			{
				Verifier: batchProof.Verifier2,
				Proof:    batchProof.BatchProof2,
			},
		}

		input, err := encoding.EncodeProveBatchesInput(metas, transitions)
		if err != nil {
			return nil, err
		}
		encodedSubProofs, err := encoding.EncodeBatchesSubProofsPacaya(subProofs)
		if err != nil {
			return nil, err
		}

		if a.proverSetAddress != rpc.ZeroAddress {
			if data, err = encoding.ProverSetPacayaABI.Pack("proveBatches", input, encodedSubProofs); err != nil {
				return nil, encoding.TryParsingCustomError(err)
			}
			to = a.proverSetAddress
		} else {
			if data, err = encoding.TaikoInboxABI.Pack("proveBatches", input, encodedSubProofs); err != nil {
				return nil, encoding.TryParsingCustomError(err)
			}
			to = a.pacayaInboxAddress
		}

		return &txmgr.TxCandidate{
			TxData:   data,
			To:       &to,
			Blobs:    nil,
			GasLimit: txOpts.GasLimit,
			Value:    txOpts.Value,
		}, nil
	}
}

// BuildProveBatchesShasta creates a new Shasta Inbox.prove transaction.
func (a *ProveBatchesTxBuilder) BuildProveBatchesShasta(
	ctx context.Context,
	batchProof *proofProducer.BatchProofs,
) TxBuilder {
	return func(txOpts *bind.TransactOpts) (*txmgr.TxCandidate, error) {
		var (
			proposals = make([]*surgeBindings.SurgeInboxClientProposed, len(batchProof.ProofResponses))
			input     = &surgeBindings.IInboxProveInput{
				Commitment: surgeBindings.IInboxCommitment{ActualProver: txOpts.From},
			}
		)

		if len(batchProof.ProofResponses) == 0 {
			return nil, fmt.Errorf("no proof responses in batch proof")
		}

		// Query contract state - needed for genesis Shasta case where we use lastFinalizedBlockHash
		// as the firstProposalParentBlockHash when there are no previous Pacaya blocks.
		var coreState *surgeBindings.IInboxCoreState
		coreState, coreStateErr := a.rpc.GetCoreStateShasta(&bind.CallOpts{Context: ctx})
		if coreStateErr != nil {
			log.Warn("Failed to get Shasta core state", "error", coreStateErr)
		} else {
			log.Debug(
				"Contract CoreState before proof submission",
				"nextProposalId", coreState.NextProposalId,
				"lastFinalizedProposalId", coreState.LastFinalizedProposalId,
				"lastFinalizedBlockHash", common.Bytes2Hex(coreState.LastFinalizedBlockHash[:]),
			)
		}

		for i, proofResponse := range batchProof.ProofResponses {
			if len(proofResponse.Opts.ShastaOptions().Headers) == 0 {
				return nil, fmt.Errorf(
					"no headers in proof response options for proposal ID %d",
					proofResponse.Meta.Shasta().GetEventData().Id,
				)
			}
			proposals[i] = proofResponse.Meta.Shasta().GetEventData()
			lastHeader := proofResponse.Opts.ShastaOptions().Headers[len(proofResponse.Opts.ShastaOptions().Headers)-1]

			proposalHash, err := a.rpc.GetShastaProposalHash(nil, proposals[i].Id)
			if err != nil {
				return nil, encoding.TryParsingCustomError(err)
			}

			// Set first proposal information.
			if i == 0 {
				input.Commitment.FirstProposalId = proposals[i].Id
				if proposals[i].Id.Cmp(common.Big1) == 0 {
					if coreState == nil {
						return nil, fmt.Errorf(
							"cannot determine firstProposalParentBlockHash for proposalId=1: coreState is nil",
						)
					}
					input.Commitment.FirstProposalParentBlockHash = coreState.LastFinalizedBlockHash
				} else {
					lastOriginInLastProposal, err := a.rpc.LastL1OriginInBatchShasta(
						ctx,
						new(big.Int).Sub(proposals[i].Id, common.Big1),
					)
					if err != nil {
						return nil, err
					}
					input.Commitment.FirstProposalParentBlockHash = lastOriginInLastProposal.L2BlockHash
				}
			}

			// Set last proposal information.
			if i == len(batchProof.ProofResponses)-1 {
				input.Commitment.LastProposalHash = proposalHash
				input.Commitment.EndBlockNumber = lastHeader.Number
				input.Commitment.EndStateRoot = lastHeader.Root
			}

			// Set transition information.
			input.Commitment.Transitions = append(input.Commitment.Transitions, surgeBindings.IInboxTransition{
				Proposer:  proposals[i].Proposer,
				Timestamp: new(big.Int).SetUint64(proofResponse.Meta.Shasta().GetTimestamp()),
				BlockHash: lastHeader.Hash(),
			})

			log.Info(
				"Build batch proof submission transaction",
				"batchID", proposals[i].Id,
				"proposalHash", proposalHash,
				"start", proofResponse.Opts.ShastaOptions().Headers[0].Number,
				"end", proofResponse.Opts.ShastaOptions().Headers[len(proofResponse.Opts.ShastaOptions().Headers)-1].Number,
				"designatedProver", batchProof.ProofResponses[i].Opts.ShastaOptions().DesignatedProver,
				"actualProver", txOpts.From,
				"firstProposalParentBlockHash", common.Bytes2Hex(input.Commitment.FirstProposalParentBlockHash[:]),
			)
		}

		// Validate consecutive proposals
		for i := 1; i < len(proposals); i++ {
			if proposals[i].Id.Uint64() != proposals[i-1].Id.Uint64()+1 {
				return nil, fmt.Errorf(
					"non-consecutive proposals: %d -> %d",
					proposals[i-1].Id,
					proposals[i].Id)
			}
		}

		inputData, err := a.rpc.EncodeProveInput(&bind.CallOpts{Context: txOpts.Context}, input)
		if err != nil {
			return nil, encoding.TryParsingCustomError(err)
		}
		log.Info(
			"Verifier information",
			"ProofType1", batchProof.ProofType1,
			"VerifierID1", batchProof.VerifierID1,
			"Proof1", common.Bytes2Hex(batchProof.BatchProof1),
			"ProofType2", batchProof.ProofType2,
			"VerifierID2", batchProof.VerifierID2,
			"Proof2", common.Bytes2Hex(batchProof.BatchProof2),
		)

		subProofs := []encoding.SubProofShasta{
			{
				VerifierId: batchProof.VerifierID1,
				Proof:      batchProof.BatchProof1,
			},
			{
				VerifierId: batchProof.VerifierID2,
				Proof:      batchProof.BatchProof2,
			},
		}
		encodedSubProofs, err := encoding.EncodeBatchesSubProofsShasta(subProofs)
		if err != nil {
			return nil, err
		}

		data, err := encoding.ShastaInboxABI.Pack("prove", inputData, encodedSubProofs)
		if err != nil {
			return nil, err
		}
		return &txmgr.TxCandidate{
			TxData:   data,
			To:       &a.shastaInboxAddress,
			Blobs:    nil,
			GasLimit: txOpts.GasLimit,
		}, nil
	}
}
