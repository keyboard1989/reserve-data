package exchange

type LiquiEndpoint interface {
	GetInfo(key string, signer Signer) (liqinfo, error)
	Depth(tokens string) (liqresp, error)
}
