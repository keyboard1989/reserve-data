package blockchain

import (
	"fmt"
	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

type Blockchain struct {
	ethclient *ethclient.Client
	wrapper   *ContractWrapper
	reserve   *ReserveContract
	signer    Signer
	tokens    []common.Token
}

func (self *Blockchain) AddToken(t common.Token) {
	self.tokens = append(self.tokens, t)
}

func (self *Blockchain) FetchBalanceData(reserve ethereum.Address) (map[string]common.BalanceEntry, error) {
	result := map[string]common.BalanceEntry{}
	tokens := []ethereum.Address{}
	for _, tok := range self.tokens {
		tokens = append(tokens, ethereum.HexToAddress(tok.Address))
	}
	// fmt.Printf("reserve: %v\n", reserve)
	// fmt.Printf("wrapper: %v\n", self.wrapper.ContractWrapperCaller.contract)
	timestamp := common.GetTimestamp()
	balances, err := self.wrapper.GetBalances(nil, reserve, tokens)
	returnTime := common.GetTimestamp()
	// fmt.Printf("balances: %v\n", balances)
	// fmt.Printf("errors: %v\n", err)
	if err != nil {
		for tokenID, _ := range common.SupportedTokens {
			result[tokenID] = common.BalanceEntry{
				Valid:      false,
				Timestamp:  timestamp,
				ReturnTime: returnTime,
			}
		}
	} else {
		for i, tok := range self.tokens {
			result[tok.ID] = common.BalanceEntry{
				Valid:      true,
				Timestamp:  timestamp,
				ReturnTime: returnTime,
				Balance:    common.RawBalance(*balances[i]),
			}
		}
	}
	return result, nil
}

func (self *Blockchain) SetRate(
	sources []ethereum.Address,
	dests []ethereum.Address,
	rates []*big.Int,
	expiryBlocks []*big.Int) (ethereum.Hash, error) {

	tx, err := self.reserve.SetRate(
		self.signer.GetTransactOpts(),
		sources, dests, rates, expiryBlocks, true)
	if err != nil {
		return ethereum.Hash{}, err
	} else {
		return tx.Hash(), err
	}
}

func (self *Blockchain) Send(
	token common.Token,
	amount *big.Int,
	dest ethereum.Address) (ethereum.Hash, error) {

	tx, err := self.reserve.Withdraw(
		self.signer.GetTransactOpts(),
		ethereum.HexToAddress(token.Address),
		amount, dest)
	if err != nil {
		return ethereum.Hash{}, err
	} else {
		return tx.Hash(), err
	}
}

// func (self *Blockchain) sendToken(token common.Token, amount *big.Int, address ethereum.Address) (ethereum.Hash, error) {
// 	erc20, err := NewErc20Contract(
// 		ethereum.HexToAddress(token.Address),
// 		self.ethclient,
// 	)
// 	fmt.Printf("address: %s\n", token.Address)
// 	if err != nil {
// 		return ethereum.Hash{}, err
// 	}
// 	tx, err := erc20.Transfer(
// 		self.signer.GetTransactOpts(),
// 		address, amount)
// 	if err != nil {
// 		return ethereum.Hash{}, err
// 	} else {
// 		return tx.Hash(), nil
// 	}
// }
//
// func (self *Blockchain) sendETH(
// 	amount *big.Int,
// 	address ethereum.Address) (ethereum.Hash, error) {
// 	// nonce, gasLimit, gasPrice gets from ethclient
//
// 	option := context.Background()
// 	rm := self.signer.GetAddress()
// 	nonce, err := self.ethclient.PendingNonceAt(
// 		option, rm)
// 	if err != nil {
// 		return ethereum.Hash{}, err
// 	}
// 	gasLimit := big.NewInt(1000000)
// 	gasPrice := big.NewInt(20000000000)
// 	rawTx := types.NewTransaction(
// 		nonce, address, amount, gasLimit, gasPrice, []byte{})
// 	signedTx, err := self.signer.Sign(rm, rawTx)
// 	if err != nil {
// 		return ethereum.Hash{}, err
// 	}
// 	if err = self.ethclient.SendTransaction(option, signedTx); err != nil {
// 		return ethereum.Hash{}, err
// 	}
// 	return signedTx.Hash(), nil
// }

func NewBlockchain(wrapperAddr, reserveAddr ethereum.Address, signer Signer) (*Blockchain, error) {
	endpoint := "http://localhost:8545"
	infura, err := ethclient.Dial(endpoint)
	if err != nil {
		return nil, err
	}
	wrapper, err := NewContractWrapper(wrapperAddr, infura)
	if err != nil {
		return nil, err
	}
	fmt.Printf("reserve address: %s\n", reserveAddr.Hex())
	reserve, err := NewReserveContract(reserveAddr, infura)
	if err != nil {
		return nil, err
	}
	return &Blockchain{
		ethclient: infura,
		wrapper:   wrapper,
		reserve:   reserve,
		signer:    signer,
		tokens:    []common.Token{},
	}, nil
}
