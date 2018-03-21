package common

import (
	"github.com/KyberNetwork/reserve-data/signer"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type HuobiHTTPConfig struct {
	Port     int
	AuthEnbl bool
	KyberENV string
}

type HuobiConfig struct {
	IntermediatorAddress ethereum.Address
	EthEndPoint          string
	StoragePath          string
	IntermediatorSigner  signer.FileSigner
	HuobihttpConfig      HuobiHTTPConfig
}
