package settings

import (
	"github.com/KyberNetwork/reserve-data/common"
)

type ExchangeStorage interface {
	// GetFee returns a map[tokenID]exchangeFees and error if occur
	GetFee(ex ExchangeName) (common.ExchangeFees, error)
	// StoreFee stores the fee with exchangeName as key into database and return error if occur
	StoreFee(ex ExchangeName, data common.ExchangeFees) error
	// GetMinDeposit returns a map[tokenID]MinDeposit and error if occur
	GetMinDeposit(ex ExchangeName) (common.ExchangesMinDeposit, error)
	// StoreMinDeposit stores the minDeposit with exchangeName as key into database and return error if occur
	StoreMinDeposit(ex ExchangeName, minDeposit common.ExchangesMinDeposit) error
	// GetDepositAddress returns a map[tokenID]DepositAddress and error if occur
	GetDepositAddress(ex ExchangeName) (common.ExchangeAddresses, error)
	// StoreDepositAddress stores the depositAddress with exchangeName as key into database and
	// return error if occur
	StoreDepositAddress(ex ExchangeName, addrs common.ExchangeAddresses) error
	// CountFee return the number of element in fee database, and error if occcur.
	// It is used mainly to check if the fee database is empty
	CountFee() (uint64, error)
}
