package alpha

import (
	"github.com/KyberNetwork/reserve-data/common"
)

type ReserveData struct {
	storage Storage
	fetcher Fetcher
}

func (self ReserveData) CurrentPriceVersion() (common.Version, error) {
	return self.storage.CurrentPriceVersion()
}

func (self ReserveData) GetAllPrices() (common.AllPriceResponse, error) {
	timestamp := common.GetTimestamp()
	version, err := self.storage.CurrentPriceVersion()
	if err != nil {
		return common.AllPriceResponse{}, err
	} else {
		result := common.AllPriceResponse{}
		data, err := self.storage.GetAllPrices(version)
		returnTime := common.GetTimestamp()
		result.Version = version
		result.Timestamp = timestamp
		result.ReturnTime = returnTime
		result.Data = data
		return result, err
	}
}

func (self ReserveData) GetOnePrice(pairID common.TokenPairID) (common.OnePriceResponse, error) {
	timestamp := common.GetTimestamp()
	version, err := self.storage.CurrentPriceVersion()
	if err != nil {
		return common.OnePriceResponse{}, err
	} else {
		result := common.OnePriceResponse{}
		data, err := self.storage.GetOnePrice(pairID, version)
		returnTime := common.GetTimestamp()
		result.Version = version
		result.Timestamp = timestamp
		result.ReturnTime = returnTime
		result.Data = data
		return result, err
	}
}

func (self ReserveData) Run() error {
	return self.fetcher.Run()
}

func NewReserveData(storage Storage, fetcher Fetcher) *ReserveData {
	return &ReserveData{storage, fetcher}
}
