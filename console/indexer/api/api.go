package api

import (
	"context"
	"energi-challenge/console/indexer/entities"
	"energi-challenge/console/indexer/utils"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

// indexer
func (idx *Indexer) indexer() {
	idx.wg.Add(idx.conf.Indexer.Workers)

	for i := 0; i < idx.conf.Indexer.Workers; i++ {
		go func() {
			defer idx.wg.Done()

			for {
				select {
				case <-idx.ctx.Done():
					return
				case i := <-idx.jobs:
					if err := idx.scan(i); err != nil {
						log.Println(err)
						continue
					}

					// check if current job newer than current latestBlock block
					if !idx.subscribed && i >= idx.latestBlock {
						idx.subscribe()

						idx.latestBlock = i

						idx.subscribed = true
					}

					log.Printf("scanned block number %d", i)
				}
			}
		}()
	}
}

// subscribe - subscribes to notifications about the current blockchain head on the given channel.
func (idx *Indexer) subscribe() {
	sub, err := idx.client.SubscribeNewHead(context.Background(), idx.events)
	if err != nil {
		log.Println(err)
		idx.subscribed = false
		return
	}

	idx.wg.Add(1)
	go func() {
		defer idx.wg.Done()

		//  log
		log.Println("subscribed to new block")
		for {
			select {
			case err := <-sub.Err():
				log.Println(err)
				idx.subscribed = false
				return
			case <-idx.ctx.Done():
				return
			case header := <-idx.events:
				head := header.Number.Int64()
				idx.jobs <- head
				idx.latestBlock = head
			}
		}
	}()
}

// scan
func (idx *Indexer) scan(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), idx.conf.Indexer.Timeout)
	defer cancel()

	// already scanned
	if idx.repository.HasScanned(ctx, id) {
		return fmt.Errorf("block %d already scanned", id)
	}

	// rate limit
	idx.limiter.Take()

	return idx.storeBlock(ctx, id)
}

// storeBlock stores the given block info into the database
func (idx *Indexer) storeBlock(ctx context.Context, id int64) error {
	// get block
	block, err := idx.client.BlockByNumber(ctx, big.NewInt(id))
	if err != nil {
		return err
	}

	hash := block.Hash()

	txCount, err := idx.client.TransactionCount(ctx, hash)
	if err != nil {
		return err
	}

	b := &entities.Block{
		Number:    id,
		Hash:      hash.Hex(),
		Timestamp: time.Unix(int64(block.Time()), 0),
		TxCount:   txCount,
	}

	if err := idx.repository.SaveBlock(ctx, b); err != nil {
		return err
	}

	// get chain id
	chainID, err := idx.client.ChainID(ctx)
	if err != nil {
		return err
	}

	for order, tx := range block.Transactions() {
		msg, err := tx.AsMessage(types.LatestSignerForChainID(chainID), nil)
		if err != nil {
			return err
		}

		t := &entities.Tx{
			BlockNumber: block.Number().Int64(),
			Hash:        tx.Hash().Hex(),
			From:        msg.From().Hex(),
			To:          utils.AddrToHex(msg.To()),
			Amount:      tx.Value().Int64(),
			Nonce:       tx.Nonce(),
			Timestamp:   time.Unix(int64(block.Time()), 0),
			Order:       order,
		}

		if err := idx.repository.SaveTx(ctx, t); err != nil {
			log.Println("idx_store_save_tx", err)
			continue
		}
	}

	return nil
}
