package market

func GetCentralizedPrice(base, quote string) (PriceData, error) {
	version := StorageInstance.GetVersion()
	timestamp := StorageInstance.GetTimestamp(version)
	data, err := StorageInstance.GetOnePairData(base, quote, version)
	if err != nil {
		return PriceData{}, err
	} else {
		return PriceData{
			Version:      version,
			Timestamp:    timestamp,
			ExchangeData: data,
		}, nil
	}
}

func GetAllCentralizedPrice() (AllPriceData, error) {
	version := StorageInstance.GetVersion()
	data, err := StorageInstance.GetAllPairData(version)
	if err != nil {
		return AllPriceData{}, err
	} else {
		return data, nil
	}
}
