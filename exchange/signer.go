package exchange

type Signer interface {
	GetLiquiKey() string
	LiquiSign(msg string) string
}
