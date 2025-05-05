package parser

import (
	"math/big"
	"strings"

	"github.com/darcygail/ether-explorer/schema"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var erc721TransferSig = "Transfer(address,address,uint256)"
var erc721SigHash = common.HexToHash("0xddf252ad...") // 事件签名hash

func ERC721TransferSigHash() common.Hash {
	return erc721SigHash
}

func IsERC721Transfer(log *types.Log) bool {
	return len(log.Topics) == 4 && log.Topics[0] == erc721SigHash
}

func ParseERC721Transfer(log *types.Log) (from, to, tokenId, contract string) {
	from = common.HexToAddress(log.Topics[1].Hex()).Hex()
	to = common.HexToAddress(log.Topics[2].Hex()).Hex()
	tokenIdBig := new(big.Int).SetBytes(log.Topics[3].Bytes())
	tokenId = tokenIdBig.String()
	contract = log.Address.Hex()
	return
}

func AddNFT(acc *schema.Account, contract string, tokenId string) {
	contract = strings.ToLower(contract)
	found := false

	for i, asset := range acc.Assets {
		if asset.Type == "ERC721" && asset.ContractAddress == contract {
			acc.Assets[i].TokenIDs = append(asset.TokenIDs, schema.TokenIDAmount{ID: tokenId})
			found = true
			break
		}
	}

	if !found {
		acc.Assets = append(acc.Assets, schema.Asset{
			ContractAddress: contract,
			Type:            "ERC721",
			TokenIDs:        []schema.TokenIDAmount{{ID: tokenId}},
		})
	}
}
