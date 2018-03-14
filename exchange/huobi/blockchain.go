package huobi

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	originalbc "github.com/KyberNetwork/reserve-data/blockchain"
	"github.com/KyberNetwork/reserve-data/blockchain/nonce"
	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/exchange"
	"github.com/KyberNetwork/reserve-data/signer"
	ether "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type tbindex struct {
	BulkIndex   uint64
	IndexInBulk uint64
}

type txExtraInfo struct {
	BlockNumber *string
	BlockHash   ethereum.Hash
	From        ethereum.Address
}

type rpcTransaction struct {
	tx *types.Transaction
	txExtraInfo
}

type Blockchain struct {
	rpcClient          *rpc.Client
	client             *ethclient.Client
	wrapper            *originalbc.KNWrapperContract
	intermediateSigner exchange.Signer
	nonceIntermediate  exchange.NonceCorpus
	chainType          string
}

func (tx *rpcTransaction) UnmarshalJSON(msg []byte) error {
	if err := json.Unmarshal(msg, &tx.tx); err != nil {
		return err
	}
	return json.Unmarshal(msg, &tx.txExtraInfo)
}

func getNextNonce(n exchange.NonceCorpus) (*big.Int, error) {
	var nonce *big.Int
	var err error
	for i := 0; i < 3; i++ {
		nonce, err = n.GetNextNonce()
		if err == nil {
			return nonce, nil
		}
	}
	return nonce, err
}

func (tx *rpcTransaction) BlockNumber() *big.Int {
	if tx.txExtraInfo.BlockNumber == nil {
		return big.NewInt(0)
	} else {
		blockno, err := hexutil.DecodeBig(*tx.txExtraInfo.BlockNumber)
		if err != nil {
			log.Printf("Error decoding block number: %v", err)
			return big.NewInt(0)
		} else {
			return blockno
		}
	}
}

func ensureContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.TODO()
	}
	return ctx
}

func donothing() {}

func (self *Blockchain) getIntermediateTransactOpts(nonce, gasPrice *big.Int) (*bind.TransactOpts, context.CancelFunc, error) {
	shared := self.intermediateSigner.GetTransactOpts()
	var err error
	if nonce == nil {
		nonce, err = getNextNonce(self.nonceIntermediate)
	}
	if err != nil {
		return nil, donothing, err
	}
	if gasPrice == nil {
		gasPrice = big.NewInt(50100000000)
	}
	timeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	result := bind.TransactOpts{
		shared.From,
		nonce,
		shared.Signer,
		shared.Value,
		gasPrice,
		shared.GasLimit,
		timeout,
	}
	return &result, cancel, nil
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
	log.Printf("amount is %d", amount)
	log.Printf("exchange address is %s", exchangeAddress.Hex())
	data, err := packData("transfer", exchangeAddress, amount)
	if err != nil {
		return nil, err
	}
	log.Printf("opts address is: %s ", opts.From.Hex())
	log.Printf("token address is: %s ", tokenAddress.Hex())
	msg := ether.CallMsg{From: opts.From, To: &tokenAddress, Value: big.NewInt(0), Data: data}
	log.Printf("message is :%x", msg)
	gasLimit, err := self.client.EstimateGas(ensureContext(opts.Context), msg)
	if err != nil {
		log.Printf("Intermediator: Can not estimate gas: %v", err)
		gasLimit = big.NewInt(25000)
	} else {
		log.Println("Intermediator: gas limit estimated is : %d", gasLimit)
		if (gasLimit.Cmp(big.NewInt(25000))) < 0 {
			gasLimit = big.NewInt(25000)
		}
	}
	//build tx, sign and send

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
		gasLimit = big.NewInt(25000)
	} else {
		log.Println("Intermediator: gas limit estimated is : %d", gasLimit)
		if (gasLimit.Cmp(big.NewInt(25000))) < 0 {
			gasLimit = big.NewInt(25000)
		}
	}
	//build tx, sign and send
	tx := types.NewTransaction(opts.Nonce.Uint64(), exchangeAddress, amount, gasLimit, opts.GasPrice, nil)
	signTX, err := self.intermediateSigner.Sign(tx)
	if err != nil {
		log.Println("Intermediator: Can not sign the transaction")
		return nil, err
	}
	log.Printf("Intermediator: the signed TX is: \n%v ", signTX)
	err = self.client.SendTransaction(ctx, signTX)
	if err != nil {
		log.Println("ERROR: Can't send the transaction")
		return nil, err
	}
	return signTX, nil
}

func (self *Blockchain) TransactionByHash(ctx context.Context, hash ethereum.Hash) (tx *rpcTransaction, isPending bool, err error) {
	var json *rpcTransaction
	log.Printf("hash is %s", hash.Hex())
	err = self.rpcClient.CallContext(ctx, &json, "eth_getTransactionByHash", hash)
	if err != nil {
		return nil, false, err
	} else if json == nil {
		return nil, false, ether.NotFound
	} else if _, r, _ := json.tx.RawSignatureValues(); r == nil {
		return nil, false, fmt.Errorf("server returned transaction without signature")
	}
	return json, json.BlockNumber == nil, nil
}

func (self *Blockchain) TxStatus(hash ethereum.Hash) (string, uint64, error) {
	option := context.Background()
	tx, pending, err := self.TransactionByHash(option, hash)
	if err == nil {
		// tx exist
		if pending {
			return "", 0, nil
		} else {
			receipt, err := self.client.TransactionReceipt(option, hash)
			if err != nil {
				log.Println("Get receipt err: ", err.Error())
				log.Printf("Receipt: %+v", receipt)
				if receipt != nil {
					// only byzantium has status field at the moment
					// mainnet, ropsten are byzantium, other chains such as
					// devchain, kovan are not
					if self.chainType == "byzantium" {
						if receipt.Status == 1 {
							// successful tx
							return "mined", tx.BlockNumber().Uint64(), nil
						} else {
							// failed tx
							return "failed", tx.BlockNumber().Uint64(), nil
						}
					} else {
						return "mined", tx.BlockNumber().Uint64(), nil
					}
				} else {
					// networking issue
					return "", 0, err
				}
			} else {
				if receipt.Status == 1 {
					// successful tx
					return "mined", tx.BlockNumber().Uint64(), nil
				} else {
					// failed tx
					return "failed", tx.BlockNumber().Uint64(), nil
				}
			}
		}
	} else {
		if err == ether.NotFound {
			// tx doesn't exist. it failed
			return "lost", 0, nil
		} else {
			// networking issue
			return "", 0, err
		}
	}
}

func NewBlockchain(intermediateSigner *signer.FileSigner, ethEndpoint string, wrapperAddr ethereum.Address) (*Blockchain, error) {
	log.Printf("intermediate address: %s", intermediateSigner.GetAddress().Hex())
	//set client & endpoint
	client, err := rpc.Dial(ethEndpoint)
	if err != nil {
		log.Println("ERROR: can't dial to etherum Endpoint ")
	}
	infura := ethclient.NewClient(client)
	intermediatenonce := nonce.NewTimeWindow(infura, intermediateSigner)
	wrapper, err := originalbc.NewKNWrapperContract(wrapperAddr, infura)
	if err != nil {
		return nil, err
	}
	return &Blockchain{
		rpcClient:          client,
		client:             infura,
		wrapper:            wrapper,
		intermediateSigner: intermediateSigner,
		nonceIntermediate:  intermediatenonce,
	}, nil
}

func (self *Blockchain) CheckBalance(token common.Token) *big.Int {
	opts, _, _ := self.getIntermediateTransactOpts(nil, nil)
	balance, err := self.FetchBalanceData(opts.From, token)
	if err != nil || !balance.Valid {
		return big.NewInt(0)
	}

	balanceFloat := balance.Balance.ToFloat(token.Decimal)
	return (getBigIntFromFloat(balanceFloat, token.Decimal))

}

func (self *Blockchain) FetchBalanceData(reserve ethereum.Address, token common.Token) (common.BalanceEntry, error) {
	result := common.BalanceEntry{}
	tokens := []ethereum.Address{}
	tokens = append(tokens, ethereum.HexToAddress(token.Address))

	timestamp := common.GetTimestamp()
	balances, err := self.wrapper.GetBalances(nil, nil, reserve, tokens)
	returnTime := common.GetTimestamp()
	log.Printf("Fetcher ------> balances: %v, err: %s", balances, err)
	if err != nil {
		result = common.BalanceEntry{
			Valid:      false,
			Error:      err.Error(),
			Timestamp:  timestamp,
			ReturnTime: returnTime,
		}
	} else {
		if balances[0].Cmp(big.NewInt(0)) == 0 || balances[0].Cmp(big.NewInt(10).Exp(big.NewInt(10), big.NewInt(33), nil)) > 0 {
			log.Printf("Fetcher ------> balances of token %s is invalid", token.ID)
			result = common.BalanceEntry{
				Valid:      false,
				Error:      "Got strange balances from node. It equals to 0 or is bigger than 10^33",
				Timestamp:  timestamp,
				ReturnTime: returnTime,
				Balance:    common.RawBalance(*balances[0]),
			}
		} else {
			result = common.BalanceEntry{
				Valid:      true,
				Timestamp:  timestamp,
				ReturnTime: returnTime,
				Balance:    common.RawBalance(*balances[0]),
			}
		}
	}

	return result, nil
}
