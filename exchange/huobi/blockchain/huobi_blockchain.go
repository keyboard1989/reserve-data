package blockchain

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"time"

	"github.com/KyberNetwork/reserve-data/blockchain/nonce"
	"github.com/KyberNetwork/reserve-data/common/blockchain"
	ether "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type Blockchain struct {
	*blockchain.BaseBlockchain
}

func packData(method string, params ...interface{}) ([]byte, error) {
	file, err := os.Open(
		"/go/src/github.com/KyberNetwork/reserve-data/blockchain/ERC20.abi")
	if err != nil {
		return nil, err
	}
	packabi, err := abi.JSON(file)
	if err != nil {
		return nil, err
	}
	data, err := packabi.Pack(method, params...)
	if err != nil {
		log.Println("Intermediator: Can not pack the data")
		return nil, err
	}
	return data, nil
}

func (self *Blockchain) SendTokenFromAccountToExchange(amount *big.Int, exchangeAddress ethereum.Address, tokenAddress ethereum.Address) (*types.Transaction, error) {
	opts, cancel, err := self.getIntermediateTransactOpts(nil, nil)
	ctx := opts.Context
	defer cancel()
	data, err := packData("transfer", exchangeAddress, amount)
	if err != nil {
		return nil, err
	}
	msg := ether.CallMsg{From: opts.From, To: &tokenAddress, Value: big.NewInt(0), Data: data}
	gasLimit, err := self.client.EstimateGas(ensureContext(opts.Context), msg)
	if err != nil {
		log.Printf("Intermediator: Can not estimate gas: %v", err)
		return nil, err
	}

	tx := types.NewTransaction(opts.Nonce.Uint64(), tokenAddress, big.NewInt(0), gasLimit, opts.GasPrice, data)
	signTX, err := self.intermediateSigner.Sign(tx)
	if err != nil {
		log.Println("Intermediator: Can not sign the transaction")
		return nil, err
	}
	err = self.client.SendTransaction(ctx, signTX)
	if err != nil {
		log.Println("Intermediator: ERROR: Can't send the transaction")
		return nil, err
	}
	return signTX, nil
}

func (self *Blockchain) SendETHFromAccountToExchange(amount *big.Int, exchangeAddress ethereum.Address) (*types.Transaction, error) {

	opts, cancel, err := self.getIntermediateTransactOpts(nil, nil)
	ctx := opts.Context
	defer cancel()
	//build msg and get gas limit
	msg := ether.CallMsg{From: opts.From, To: &exchangeAddress, Value: amount, Data: nil}
	gasLimit, err := self.client.EstimateGas(ensureContext(opts.Context), msg)
	if err != nil {
		log.Printf("Intermediator: Can not estimate gas: %v", err)
		return nil, err
	}
	//build tx, sign and send
	tx := types.NewTransaction(opts.Nonce.Uint64(), exchangeAddress, amount, gasLimit, opts.GasPrice, nil)
	signTX, err := self.intermediateSigner.Sign(tx)
	if err != nil {
		log.Println("Intermediator: Can not sign the transaction")
		return nil, err
	}
	err = self.client.SendTransaction(ctx, signTX)
	if err != nil {
		log.Println("ERROR: Can't send the transaction")
		return nil, err
	}
	return signTX, nil
}

func NewBlockchain(intermediateSigner Signer, ethEndpoint string) (*Blockchain, error) {
	log.Printf("intermediate address: %s", intermediateSigner.GetAddress().Hex())
	//set client & endpoint
	client, err := rpc.Dial(ethEndpoint)
	if err != nil {
		log.Println("ERROR: can't dial to etherum Endpoint ")
	}
	infura := ethclient.NewClient(client)
	intermediatenonce := nonce.NewTimeWindow(infura, intermediateSigner)
	// wrp, err := originalbc.NewKNWrapperContract(wrapperAddr, infura)
	if err != nil {
		return nil, err
	}
	return &Blockchain{
		rpcClient: client,
		client:    infura,
		// wrapper:            wrp,
		intermediateSigner: intermediateSigner,
		nonceIntermediate:  intermediatenonce,
	}, nil
}

// func (self *Blockchain) CheckBalance(token common.Token) *big.Int {
// 	addr := self.intermediateSigner.GetAddress()
// 	balance, err := self.FetchBalanceData(addr, token)
// 	if err != nil || !balance.Valid {
// 		return big.NewInt(0)
// 	}

// 	balanceFloat := balance.Balance.ToFloat(token.Decimal)
// 	return (getBigIntFromFloat(balanceFloat, token.Decimal))

// }

// func (self *Blockchain) FetchBalanceData(reserve ethereum.Address, token common.Token) (common.BalanceEntry, error) {
// 	result := common.BalanceEntry{}
// 	tokens := []ethereum.Address{}
// 	tokens = append(tokens, ethereum.HexToAddress(token.Address))

// 	timestamp := common.GetTimestamp()
// 	balances, err := self.wrapper.GetBalances(nil, nil, reserve, tokens)
// 	returnTime := common.GetTimestamp()
// 	log.Printf("Fetcher ------> balances: %v, err: %s", balances, err)
// 	if err != nil {
// 		result = common.BalanceEntry{
// 			Valid:      false,
// 			Error:      err.Error(),
// 			Timestamp:  timestamp,
// 			ReturnTime: returnTime,
// 		}
// 	} else {
// 		if balances[0].Cmp(big.NewInt(0)) == 0 || balances[0].Cmp(big.NewInt(10).Exp(big.NewInt(10), big.NewInt(33), nil)) > 0 {
// 			log.Printf("Fetcher ------> balances of token %s is invalid", token.ID)
// 			result = common.BalanceEntry{
// 				Valid:      false,
// 				Error:      "Got strange balances from node. It equals to 0 or is bigger than 10^33",
// 				Timestamp:  timestamp,
// 				ReturnTime: returnTime,
// 				Balance:    common.RawBalance(*balances[0]),
// 			}
// 		} else {
// 			result = common.BalanceEntry{
// 				Valid:      true,
// 				Timestamp:  timestamp,
// 				ReturnTime: returnTime,
// 				Balance:    common.RawBalance(*balances[0]),
// 			}
// 		}
// 	}

// 	return result, nil
// }
