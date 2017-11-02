package exchange

import (
	"math/big"

	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type LiquiEndpoint interface {
	GetInfo(key string, signer Signer) (liqinfo, error)
	Withdraw(string, common.Token, *big.Int, ethereum.Address, Signer) error
	Depth(tokens string) (liqresp, error)
}
