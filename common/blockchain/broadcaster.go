package blockchain

import (
	"context"
	"errors"
	"log"
	"math/big"
	"sync"
	"time"

	ether "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const DEADLINE = 1*time.Second - 10*time.Millisecond

// Broadcaster takes a signed tx and try to broadcast it to all
// nodes that it manages as fast as possible. It returns a map of
// failures and a bool indicating that the tx is broadcasted to
// at least 1 node
type Broadcaster struct {
	clients map[string]*ethclient.Client
}

func (self Broadcaster) broadcast(
	ctx context.Context,
	id string, client *ethclient.Client, tx *types.Transaction,
	wg *sync.WaitGroup, failures *sync.Map) {
	defer wg.Done()
	err := client.SendTransaction(ctx, tx)
	if err != nil {
		failures.Store(id, err)
	}
}

func (self Broadcaster) Broadcast(tx *types.Transaction) (map[string]error, bool) {
	failures := sync.Map{}
	wg := sync.WaitGroup{}
	for id, client := range self.clients {
		wg.Add(1)
		timeout, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		self.broadcast(timeout, id, client, tx, &wg, &failures)
		defer cancel()
	}
	wg.Wait()
	result := map[string]error{}
	failures.Range(func(key, value interface{}) bool {
		result[key.(string)] = value.(error)
		return true
	})
	return result, len(result) != len(self.clients) && len(self.clients) > 0
}

func NewBroadcaster(clients map[string]*ethclient.Client) *Broadcaster {
	return &Broadcaster{
		clients: clients,
	}
}

func (self Broadcaster) CallContract(ctx context.Context, msg ether.CallMsg, blockNo *big.Int, callOder []string) (output []byte, err error) {
	for _, urlstring := range callOder {
		client, _ := self.clients[urlstring]
		ctx, cancel := context.WithTimeout(ctx, DEADLINE)
		defer cancel()
		output, err = client.CallContract(ctx, msg, blockNo)
		select {
		case <-time.After(DEADLINE):
			log.Printf("FALLBACK: Ether client %s time out, trying next one...", urlstring)
		case <-ctx.Done():
			if err != nil {
				log.Printf("FALLBACK: Ether client %s done, getting err %v, trying next one...", urlstring, err)
			} else {
				log.Printf("FALLBACK: Ether client %s done, returnning result...", urlstring)
				return
			}
		}
	}
	if err == nil {
		return nil, errors.New("No eth endpoint is responsive")
	}
	return
}
