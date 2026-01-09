package types

import "math/big"

type ProposalsRollbacked struct {
	FirstProposalId *big.Int
	LastProposalId  *big.Int
}

type ProposalsRollbackedRanges []ProposalsRollbacked

// Contains checks if the given proposal ID is in the proposals rollbacked ranges.
func (r ProposalsRollbackedRanges) Contains(proposalId *big.Int) bool {
	if proposalId == nil {
		return false
	}
	for _, interval := range r {
		if interval.FirstProposalId == nil || interval.LastProposalId == nil {
			continue
		}
		if proposalId.Cmp(interval.FirstProposalId) >= 0 && proposalId.Cmp(interval.LastProposalId) <= 0 {
			return true
		}
	}
	return false
}
