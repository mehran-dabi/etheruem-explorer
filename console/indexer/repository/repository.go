package repository

import (
	"context"
	"database/sql"
	"energi-challenge/console/indexer/entities"
)

type IIndexerRepository interface {
	HasScanned(ctx context.Context, id int64) bool
	SaveBlock(ctx context.Context, block *entities.Block) error
	SaveTx(ctx context.Context, tx *entities.Tx) error
}

type IndexerRepository struct {
	db *sql.DB
}

func NewIndexerRepository(db *sql.DB) *IndexerRepository {
	return &IndexerRepository{db: db}
}

// HasScanned checks if the block has already been scanned
func (i *IndexerRepository) HasScanned(ctx context.Context, id int64) bool {
	var found = 0
	if err := i.db.QueryRowContext(
		ctx,
		hasScannedQuery,
		id,
	).Scan(&found); err != nil {
		return false
	}

	return found != 0
}

// SaveBlock saved the given block in the database
func (i *IndexerRepository) SaveBlock(ctx context.Context, b *entities.Block) error {
	_, err := i.db.ExecContext(
		ctx,
		insertBlockQuery,
		b.Number,
		b.Hash,
		b.Timestamp,
		b.TxCount,
	)
	return err
}

// SaveTx saves the given transaction into the database
func (i *IndexerRepository) SaveTx(ctx context.Context, tx *entities.Tx) error {
	_, err := i.db.ExecContext(
		ctx,
		insertTx,
		tx.Hash,
		tx.BlockNumber,
		tx.From,
		tx.To,
		tx.Amount,
		tx.Nonce,
		tx.Timestamp,
		tx.Order,
	)

	return err
}
