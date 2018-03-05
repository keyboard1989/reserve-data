package blockchain

import (
	ethereum "github.com/ethereum/go-ethereum/common"
)

type Receipt struct {
	// Consensus fields
	Status            uint   `json:"status"`
	CumulativeGasUsed uint64 `json:"cumulativeGasUsed" gencodec:"required"`

	// Implementation fields (don't reorder!)
	TxHash          ethereum.Hash    `json:"transactionHash" gencodec:"required"`
	ContractAddress ethereum.Address `json:"contractAddress"`
	GasUsed         uint64           `json:"gasUsed" gencodec:"required"`
}
