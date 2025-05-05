package worker

import (
	"context"
	"log"
	"math/big"

	"github.com/darcygail/ether-explorer/internal/eth"
	"github.com/darcygail/ether-explorer/internal/store"
	"github.com/darcygail/ether-explorer/schema"
)

func Run(workConfig schema.WorkConfig) {
	fetcher, err := eth.NewFetcher(workConfig.RpcUrl)
	if err != nil {
		log.Fatal("init fetcher failed")
		return
	}

	err = store.InitDB(workConfig.MongoUri)
	if err != nil {
		log.Fatal("init fetcher failed")
		return
	}

	ctx := context.Background()
	fetchAndSave := func(idx int64) {
		accounts, err := fetcher.Fetch(ctx, big.NewInt(idx))
		if err != nil {
			log.Printf("fetch block accounts failed, block number:%d,err:%+v", idx, err)
		}
		for _, account := range accounts {
			store.SaveAsset(ctx, account)
		}
	}

	latestBlockNumber, err := fetcher.GetLatestBlockNumber(ctx)
	bm := latestBlockNumber.Int64()

	// TODO concurrently do sync
	if !workConfig.Full {
		nBlocks := 10

		minBm := bm - int64(nBlocks)
		if err != nil {
			log.Fatal("init fetcher failed")
			return
		}
		for i := minBm; i <= bm; i++ {
			fetchAndSave(i)
		}
	} else {
		for i := int64(1); i < bm; i++ {
			fetchAndSave(i)
		}
	}
}
