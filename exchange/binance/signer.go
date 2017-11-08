package binance

type Signer interface {
	GetBinanceKey() string
	BinanceSign(msg string) string
}
