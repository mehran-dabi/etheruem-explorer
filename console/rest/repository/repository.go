package repository

import (
	"context"
	"database/sql"
	"energi-challenge/console/rest/entities"
)

type IRestRepository interface {
	GetLatestBlock(ctx context.Context) (*entities.Block, error)
	GetBlock(ctx context.Context, n int64) (*entities.Block, error)
	GetLatestTx(ctx context.Context) (*entities.Tx, error)
	GetTx(ctx context.Context, hash string) (*entities.Tx, error)
	GetStats(ctx context.Context, i, j int64) (*entities.Stats, error)
}

type RestRepository struct {
	db *sql.DB
}

func NewRestRepository(db *sql.DB) *RestRepository {
	return &RestRepository{db: db}
}

// GetLatestBlock get the latest block info from database
func (r *RestRepository) GetLatestBlock(ctx context.Context) (*entities.Block, error) {
	// get the latest block
	var b entities.Block
	if err := r.db.QueryRowContext(
		ctx,
		selectLatestBlockQuery,
	).Scan(
		&b.Number,
		&b.Hash,
		&b.Timestamp,
		&b.TxCount,
	); err != nil {
		return nil, err
	}

	// get block transaction
	rows, err := r.db.QueryContext(
		ctx,
		selectAllTxsHashesByBlockIDQuery,
		b.Number,
	)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var txs []string

	for rows.Next() {
		var tx string
		if err := rows.Scan(&tx); err != nil {
			return nil, err
		}

		txs = append(txs, tx)
	}

	b.Txs = txs

	return &b, nil
}

// GetBlock get the given block number information from database
func (r *RestRepository) GetBlock(ctx context.Context, n int64) (*entities.Block, error) {
	// get block
	var b entities.Block
	if err := r.db.QueryRowContext(
		ctx,
		selectBlockQuery,
		n,
	).Scan(
		&b.Number,
		&b.Hash,
		&b.Timestamp,
		&b.TxCount,
	); err != nil {
		return nil, err
	}

	// get block transaction
	rows, err := r.db.QueryContext(
		ctx,
		selectAllTxsHashesByBlockIDQuery,
		b.Number,
	)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var txs []string

	for rows.Next() {
		var tx string
		if err := rows.Scan(&tx); err != nil {
			return nil, err
		}

		txs = append(txs, tx)
	}

	b.Txs = txs

	return &b, nil
}

// GetLatestTx get the latest transaction information from database
func (r *RestRepository) GetLatestTx(ctx context.Context) (*entities.Tx, error) {
	var t entities.Tx
	if err := r.db.QueryRowContext(
		ctx,
		selectLatestTxQuery,
	).Scan(
		&t.Hash,
		&t.BlockNumber,
		&t.From,
		&t.To,
		&t.Amount,
		&t.Nonce,
		&t.Timestamp,
		&t.Order,
	); err != nil {
		return nil, err
	}
	return &t, nil
}

// GetTx gets the given transaction hash information from database
func (r *RestRepository) GetTx(ctx context.Context, hash string) (*entities.Tx, error) {
	var t entities.Tx
	if err := r.db.QueryRowContext(
		ctx,
		selectTxQuery,
		hash,
	).Scan(
		&t.Hash,
		&t.BlockNumber,
		&t.From,
		&t.To,
		&t.Amount,
		&t.Nonce,
		&t.Timestamp,
		&t.Order,
	); err != nil {
		return nil, err
	}
	return &t, nil
}

// GetStats gets the transaction status for the given block number range
func (r *RestRepository) GetStats(ctx context.Context, i, j int64) (*entities.Stats, error) {
	var status = &entities.Stats{
		Txs:         []string{},
		TotalAmount: 0,
	}

	rows, err := r.db.QueryContext(ctx, selectAllTxHashQuery, i, j)
	if err != nil {
		if err == sql.ErrNoRows {
			return status, nil
		}
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var tx string
		if err := rows.Scan(&tx); err != nil {
			return nil, err
		}

		status.Txs = append(status.Txs, tx)
	}

	var total float64
	if err := r.db.QueryRowContext(ctx, selectSumOfAllTxQuery, i, j).Scan(&total); err != nil {
		return nil, err
	}

	status.TotalAmount = total

	return status, nil
}
