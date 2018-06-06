package settings

import "github.com/KyberNetwork/reserve-data/common"

// GetFee returns a map[tokenID]exchangeFees and error if occur
func (setting *Settings) GetFee(ex ExchangeName) (common.ExchangeFees, error) {
	return setting.Exchange.Storage.GetFee(ex)
}

// StoreFee stores the fee with exchangeName as key into database and return error if occur
func (setting *Settings) StoreFee(ex ExchangeName, data common.ExchangeFees) error {
	return setting.Exchange.Storage.StoreFee(ex, data)
}

// GetMinDeposit returns a map[tokenID]MinDeposit and error if occur
func (setting *Settings) GetMinDeposit(ex ExchangeName) (common.ExchangesMinDeposit, error) {
	return setting.Exchange.Storage.GetMinDeposit(ex)
}

// StoreMinDeposit stores the minDeposit with exchangeName as key into database and return error if occur
func (setting *Settings) StoreMinDeposit(ex ExchangeName, minDeposit common.ExchangesMinDeposit) error {
	return setting.Exchange.Storage.StoreMinDeposit(ex, minDeposit)
}

// GetDepositAddress returns a map[tokenID]DepositAddress and error if occur
func (setting *Settings) GetDepositAddress(ex ExchangeName) (common.ExchangeAddresses, error) {
	return setting.Exchange.Storage.GetDepositAddress(ex)
}

// StoreDepositAddress stores the depositAddress with exchangeName as key into database and
// return error if occur
func (setting *Settings) StoreDepositAddress(ex ExchangeName, addrs common.ExchangeAddresses) error {
	return setting.Exchange.Storage.StoreDepositAddress(ex, addrs)
}

// GetTokenPairs returns a list of TokenPairs available at current exchange
// return error if occur
func (setting *Settings) GetTokenPairs(ex ExchangeName) ([]common.TokenPair, error) {
	return setting.Exchange.Storage.GetTokenPairs(ex)
}

// StoreTokenPairs store the list of TokenPairs with exchangeName as key into database and
// return error if occur
func (setting *Settings) StoreTokenPairs(ex ExchangeName, data []common.TokenPair) error {
	return setting.Exchange.Storage.StoreTokenPairs(ex, data)
}
