package blockchain

import (
	"context"
	"log"
	"math/big"
	"time"

	ether "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/ethclient"
)

const TIMEOUT = 2 * time.Second

type CallClients struct {
	clients []*ethclient.Client
	urls    []string
}

func NewCallClients(clients []*ethclient.Client, urls []string) *CallClients {
	return &CallClients{
		clients: clients,
		urls:    urls,
	}
}

func (self CallClients) CallContract(ctx context.Context, msg ether.CallMsg, blockNo *big.Int) (output []byte, err error) {
	for i, client := range self.clients {
		urlstring := self.urls[i]
		ctx, cancel := context.WithTimeout(ctx, TIMEOUT)
		defer cancel()
		output, err = client.CallContract(ctx, msg, blockNo)
		if err != nil {
			log.Printf("FALLBACK: Ether client %s done, getting err %v, trying next one...", urlstring, err)
		} else {
			log.Printf("FALLBACK: Ether client %s done, returnning result...", urlstring)
			return
		}

	}
	return
}
