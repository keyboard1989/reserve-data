package stat

import (
	"time"
)

type CMCEthUSDRate struct {
	cachedRates       [][]float64
	currentCacheMonth uint64
	realtimeTimepoint uint64
	realtimeRate      float64
}

func GetTimeStamp(year int, month time.Month, day int, hour int, minute int, sec int, nanosec int, loc *time.Location) uint64 {
	return uint64(time.Date(year, month, day, hour, minute, sec, nanosec, loc).Unix() * 1000)
}

func GetMonthTimeStamp(timepoint uint64) uint64 {
	t := time.Unix(int64(timepoint/1000), 0).UTC()
	month, year := t.Month(), t.Year()
	return GetTimeStamp(year, month, 1, 0, 0, 0, 0, time.UTC)
}

func (self *CMCEthUSDRate) GetUSDRate(timepoint uint64) float64 {
	if timepoint >= self.realtimeTimepoint {
		return self.realtimeRate
	}
	return self.rateFromCache(timepoint)
}

func (self *CMCEthUSDRate) rateFromCache(timepoint uint64) float64 {
	if GetMonthTimestamp(timepoint) != self.currentCacheMonth {
		// query to CMC to get a month of data
		// TODO
	} else {
		rate, err := findEthRate(self.cachedRates, timepoint)
		if err != nil {
			return self.realtimeRate
		} else {
			return rate
		}
	}
}

func (self *CMCEthUSDRate) RunGetEthRate() {
	tick := time.NewTicker(10 * time.Minute)
	go func() {
		for {
			err := self.FetchEthRate()
			if err != nil {
				log.Println(err)
			}
			<-tick.C
		}
	}()
}

func (self *CMCEthUSDRate) FetchEthRate() (err error) {
	resp, err := http.Get("https://api.coinmarketcap.com/v1/ticker/?convert=USD&limit=10")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	rateResponse := CoinCapRateResponse{}
	err = json.Unmarshal(body, &rateResponse)
	if err != nil {
		log.Printf("Getting eth-usd rate failed: %+v", err)
	} else {
		for _, rate := range rateResponse {
			if rate.Symbol == "ETH" {
				newrate, err := strconv.ParseFloat(rate.PriceUSD, 64)
				if err != nil {
					log.Println("Cannot get usd rate: %s", err.Error())
					return err
				} else {
					if self.realtimeRate == 0 {
						// set realtimeTimepoint to the timepoint that realtime rate is updated for the
						// first time
						self.realtimeTimepoint = common.GetTimepoint()
					}
					self.realtimeRate = newrate
					return nil
				}
			}
		}
	}
	return nil
}

func (self *CMCEthUSDRate) Run() {
	// run real time fetcher
	self.RunGetEthRate()
}

func NewCMCEthUSDRate() *CMCEthUSDRate {
	result := &CMCEthUSDRate{}
	result.Run()
	return result
}
