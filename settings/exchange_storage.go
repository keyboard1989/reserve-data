package settings

import (
	"github.com/KyberNetwork/reserve-data/common"
)

type ExchangeStorage interface {
	// GetFee returns a map[tokenID]exchangeFees and error if occur
	GetFee(ex ExchangeName) (map[string]common.ExchangeFees, error)
	// StoreFee stores the fee with exchangeName as key into database and return error if occur
	StoreFee(ex ExchangeName, data map[string]common.Exchange) error
	// GetMinDeposit returns a map[tokenID]MinDeposit and error if occur
	GetMinDeposit(ex ExchangeName) (common.ExchangesMinDeposit, error)
	// StoreMinDeposit stores the minDeposit with exchangeName as key into database and return error if occur
	StoreMinDeposit(ex ExchangeName, minDeposit common.ExchangesMinDeposit) error
	// GetDepositAddress returns a map[tokenID]DepositAddress and error if occur
	GetDepositAddress(ex ExchangeName) (common.ExchangeAddresses, error)
	// StoreDepositAddress stores the depositAddress with exchangeName as key into database and
	// return error if occur
	StoreDepositAddress(ex ExchangeName, addrs common.ExchangeAddresses) error
}
