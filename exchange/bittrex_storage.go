package exchange

// This storage is used to store deposit histories
// to check if a deposit history id is registered for
// a deposit already
type BittrexStorage interface {
	IsNewBittrexDeposit(id uint64) bool
	RegisterBittrexDeposit(id uint64) error
}
