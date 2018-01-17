package blockchain

import (
	"context"
	"math/big"
	"time"
)

func (self *Blockchain) InterpretTimestamp(blockno uint64, txindex uint) (uint64, error) {
	context := context.Background()
	block, err := self.client.HeaderByNumber(context, big.NewInt(int64(blockno)))
	if err != nil {
		if block == nil {
			return uint64(0), err
		} else {
			// error because parity and geth are not compatible in mix hash
			// so we ignore it
			err = nil
		}
	}
	unixSecond := block.Time.Uint64()
	unixNano := uint64(time.Unix(int64(unixSecond), 0).UnixNano())
	result := unixNano + uint64(txindex)
	return result, nil
}
