package bitfinex

type Signer interface {
	GetBitfinexKey() string
	BitfinexSign(msg string) string
}
