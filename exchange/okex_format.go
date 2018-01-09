package exchange

type OkexInfo struct {
	Data []struct {
		BaseCurrency   int     `json:"baseCurrency"`
		PricePrecision int     `json:"maxPriceDigit"`
		AmounPrecision int     `json:"maxSizeDigit"`
		MinAmount      float32 `json:"minTradeSize"`
		Symbol         string  `json:"symbol"`
	} `json:"data"`
}

type OkexAccountInfo struct {
	Result bool `json:"result"`
	Info   struct {
		Funds struct {
			Borrow  map[string]string `json:"borrow"`
			Free    map[string]string `json:"free"`
			Freezed map[string]string `json:"freezed"`
		} `json:"funds"`
	} `json:"info"`
}

type OkexDepth struct {
	Asks [][]float64 `json:"asks"`
	Bids [][]float64 `json:"bids"`
}

type OkexTrade struct {
	Result  bool   `json:"result"`
	OrderID uint64 `json:"order_id"`
}

type OkexCancel struct {
	Result  bool `json:"result"`
	OrderID int  `json:"order_id"`
}

type OkexOrderStatus struct {
	Result bool `json:"result"`
	Orders []struct {
		Amount     float32 `json:"amount"`
		AvgPrice   float32 `json:"avg_price"`
		CreateDate uint64  `json:"create_date"`
		DealAmount float32 `json:"deal_amount"`
		OrderID    int     `json:"order_id"`
		Price      float32 `json:"price"`
		Status     int     `json:"status"`
		Symbol     string  `json:"symbol"`
		Type       string  `json:"type"`
	} `json:"orders"`
}
