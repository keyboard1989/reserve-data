package fetcher

import (
	ethereum "github.com/ethereum/go-ethereum/common"
)

type Setting interface {
	GetAddress(name string) (ethereum.Address, error)
}
