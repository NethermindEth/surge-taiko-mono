package encoding

import (
	"math/big"

	"github.com/ethereum/go-ethereum/params"

	"github.com/taikoxyz/taiko-mono/packages/taiko-client/bindings"
)

const (
	SurgeNetworkID     = 763373 // 0xba5ed
	SurgeTestNetworkID = 763374 // 0xba5ee
)

var (
	livenessBond, _             = new(big.Int).SetString("125000000000000000000", 10)
	surgeLivenessBond, _        = new(big.Int).SetString("70000000000000000", 10)
	InternlDevnetProtocolConfig = &bindings.TaikoDataConfig{
		ChainId:               params.TaikoInternalL2ANetworkID.Uint64(),
		BlockMaxProposals:     324_000,
		BlockRingBufferSize:   360_000,
		MaxBlocksToVerify:     16,
		BlockMaxGasLimit:      240_000_000,
		LivenessBond:          livenessBond,
		StateRootSyncInternal: 1,
		MaxAnchorHeightOffset: 64,
		OntakeForkHeight:      1,
		BaseFeeConfig: bindings.LibSharedDataBaseFeeConfig{
			AdjustmentQuotient:     8,
			SharingPctg:            0,
			GasIssuancePerSecond:   5_000_000,
			MinGasExcess:           1_340_000_000,
			MaxGasIssuancePerBlock: 600_000_000,
		},
		MaxLivenessDisruptionPeriod: 604800,
	}
	HeklaProtocolConfig = &bindings.TaikoDataConfig{
		ChainId:               params.HeklaNetworkID.Uint64(),
		BlockMaxProposals:     324_000,
		BlockRingBufferSize:   324_512,
		MaxBlocksToVerify:     16,
		BlockMaxGasLimit:      240_000_000,
		LivenessBond:          livenessBond,
		StateRootSyncInternal: 1,
		MaxAnchorHeightOffset: 64,
		OntakeForkHeight:      1,
		BaseFeeConfig: bindings.LibSharedDataBaseFeeConfig{
			AdjustmentQuotient:     8,
			SharingPctg:            0,
			GasIssuancePerSecond:   5_000_000,
			MinGasExcess:           1_340_000_000,
			MaxGasIssuancePerBlock: 600_000_000,
		},
		MaxLivenessDisruptionPeriod: 604800,
	}
	MainnetProtocolConfig = &bindings.TaikoDataConfig{
		ChainId:               params.TaikoMainnetNetworkID.Uint64(),
		BlockMaxProposals:     324_000,
		BlockRingBufferSize:   360_000,
		MaxBlocksToVerify:     16,
		BlockMaxGasLimit:      240_000_000,
		LivenessBond:          livenessBond,
		StateRootSyncInternal: 1,
		MaxAnchorHeightOffset: 64,
		OntakeForkHeight:      1,
		BaseFeeConfig: bindings.LibSharedDataBaseFeeConfig{
			AdjustmentQuotient:     8,
			GasIssuancePerSecond:   5_000_000,
			MinGasExcess:           1_340_000_000,
			MaxGasIssuancePerBlock: 600_000_000,
		},
		MaxLivenessDisruptionPeriod: 604800,
	}
	SurgeProtocolConfig = &bindings.TaikoDataConfig{
		ChainId:               SurgeNetworkID,
		BlockMaxProposals:     324_000,
		BlockRingBufferSize:   360_000,
		MaxBlocksToVerify:     4,
		BlockMaxGasLimit:      600_000_000,
		LivenessBond:          surgeLivenessBond,
		StateRootSyncInternal: 1,
		MaxAnchorHeightOffset: 64,
		OntakeForkHeight:      1,
		BaseFeeConfig: bindings.LibSharedDataBaseFeeConfig{
			AdjustmentQuotient:     8,
			GasIssuancePerSecond:   100_000_000,
			MinGasExcess:           31_136_000_000,
			MaxGasIssuancePerBlock: 6_000_000_000,
		},
		MaxLivenessDisruptionPeriod: 604800,
	}
	SurgeTestnetProtocolConfig = &bindings.TaikoDataConfig{
		ChainId:               SurgeTestNetworkID,
		BlockMaxProposals:     324_000,
		BlockRingBufferSize:   360_000,
		MaxBlocksToVerify:     4,
		BlockMaxGasLimit:      600_000_000,
		LivenessBond:          surgeLivenessBond,
		StateRootSyncInternal: 1,
		MaxAnchorHeightOffset: 64,
		OntakeForkHeight:      1,
		BaseFeeConfig: bindings.LibSharedDataBaseFeeConfig{
			AdjustmentQuotient:     8,
			GasIssuancePerSecond:   100_000_000,
			MinGasExcess:           31_136_000_000,
			MaxGasIssuancePerBlock: 6_000_000_000,
		},
		MaxLivenessDisruptionPeriod: 604800,
	}
)

// GetProtocolConfig returns the protocol config for the given chain ID.
func GetProtocolConfig(chainID uint64) *bindings.TaikoDataConfig {
	switch chainID {
	case params.HeklaNetworkID.Uint64():
		return HeklaProtocolConfig
	case params.TaikoMainnetNetworkID.Uint64():
		return MainnetProtocolConfig
	case SurgeNetworkID:
		return SurgeProtocolConfig
	case SurgeTestNetworkID:
		return SurgeTestnetProtocolConfig
	default:
		return InternlDevnetProtocolConfig
	}
}
