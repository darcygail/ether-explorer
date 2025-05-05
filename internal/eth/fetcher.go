package eth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/darcygail/ether-explorer/internal/parser"
	"github.com/darcygail/ether-explorer/schema"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 1. fetch by block number, get all transactions and accounts from block
// type Fetcher struct{

// }
// latestBlock, _ := eth.EthClient.BlockByNumber(context.Background(), nil)

// func blockNumberCheck
var ERC721TransferEventHash = crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))

type Fetcher struct {
	etherClient *ethclient.Client
}

func (f *Fetcher) GetLatestBlockNumber(ctx context.Context) (*big.Int, error) {
	blockNumber, err := f.etherClient.BlockNumber(context.Background())
	if err != nil {
		return nil, err
	}
	return big.NewInt(int64(blockNumber)), nil
}

func (f *Fetcher) Fetch(ctx context.Context, blockNumber *big.Int) ([]*schema.Account, error) {
	accountMap := make(map[string]*schema.Account)
	//TODO check block number here or outside
	chainID, err := f.etherClient.NetworkID(context.Background())
	if err != nil {
		return nil, err
	}
	signer := types.NewLondonSigner(chainID)
	block, _ := f.etherClient.BlockByNumber(ctx, blockNumber)

	// parse all transaction
	for _, tx := range block.Transactions() {
		// parse ether accounts
		from, err := types.Sender(signer, tx)
		if err != nil {
			return nil, err
		}
		var to common.Address
		if tx.To() != nil {
			to = *tx.To()
		}

		for _, addr := range []common.Address{from, to} {
			if addr.Hex() == "" {
				continue
			}
			if _, exists := accountMap[addr.Hex()]; !exists {
				balance, _ := f.etherClient.BalanceAt(ctx, addr, blockNumber)
				accountMap[addr.Hex()] = &schema.Account{
					Address: addr.Hex(),
					Balance: balance.String(),
					Assets:  []schema.Asset{},
				}
			}
		}

		receipt, err := f.etherClient.TransactionReceipt(ctx, tx.Hash())
		if err != nil {
			log.Printf("Failed to get receipt for tx %s: %v", tx.Hash().Hex(), err)
			continue
		}

		for _, lg := range receipt.Logs {
			if parser.IsERC721Transfer(lg) {
				from, to, tokenId, contract := parser.ParseERC721Transfer(lg)

				for _, addr := range []string{from, to} {
					if _, exists := accountMap[addr]; !exists {
						balance, _ := f.etherClient.BalanceAt(ctx, common.Address(common.FromHex(addr)), blockNumber)
						accountMap[addr] = &schema.Account{
							Address: addr,
							Balance: balance.String(),
						}
					}
				}

				// 将 NFT 添加到接收方资产中
				toAcc := accountMap[to]
				parser.AddNFT(toAcc, contract, tokenId)
			}
		}

	}
	var accounts = []*schema.Account{}
	// 最终保存所有账户
	for _, acc := range accountMap {
		accounts = append(accounts, acc)
	}
	return accounts, nil
}

func (f *Fetcher) queryERC721Balance(contract common.Address, owner common.Address, block *big.Int) (*big.Int, error) {
	erc721ABI, err := abi.JSON(strings.NewReader(`[{"constant":true,"inputs":[{"name":"owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"type":"function"}]`))
	if err != nil {
		return nil, err
	}

	data, err := erc721ABI.Pack("balanceOf", owner)
	if err != nil {
		return nil, err
	}

	msg := ethereum.CallMsg{
		To:   &contract,
		Data: data,
	}

	result, err := f.etherClient.CallContract(context.Background(), msg, block)
	if err != nil {
		return nil, err
	}

	var balance *big.Int
	err = erc721ABI.UnpackIntoInterface(&balance, "balanceOf", result)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func NewFetcher(rpcUrl string) (*Fetcher, error) {
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("dial rpc failed,rpc url:%s,err:%s", rpcUrl, err))
	}
	return &Fetcher{
		etherClient: client,
	}, nil
}
