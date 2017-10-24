package blockchain

import (
	"fmt"
	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Blockchain struct {
	client *ContractWrapper
	tokens []common.Token
}

func (self *Blockchain) AddToken(t common.Token) {
	self.tokens = append(self.tokens, t)
}

func (self *Blockchain) FetchBalanceData(reserve ethereum.Address) (map[string]common.RawBalance, error) {
	result := map[string]common.RawBalance{}
	tokens := []ethereum.Address{}
	for _, tok := range self.tokens {
		tokens = append(tokens, ethereum.HexToAddress(tok.Address))
	}
	fmt.Printf("reserve: %v\n", reserve)
	fmt.Printf("client: %v\n", self.client.ContractWrapperCaller.contract)
	balances, err := self.client.GetBalances(nil, reserve, tokens)
	fmt.Printf("balances: %v\n", balances)
	fmt.Printf("errors: %v\n", err)
	if err != nil {
		return result, err
	}
	for i, tok := range self.tokens {
		result[tok.ID] = common.RawBalance(*balances[i])
	}
	return result, nil
}

func NewBlockchain(addr ethereum.Address) (*Blockchain, error) {
	endpoint := "https://kovan.infura.io"
	infura, err := ethclient.Dial(endpoint)
	if err != nil {
		return nil, err
	}
	fmt.Printf("wrapper address: %s\n", addr.Hex())
	client, err := NewContractWrapper(addr, infura)
	if err != nil {
		return nil, err
	}
	return &Blockchain{
		client: client,
		tokens: []common.Token{},
	}, nil
}
