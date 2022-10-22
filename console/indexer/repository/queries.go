package repository

const (
	blocksTableName       = "blocks"
	transactionsTableName = "transactions"
)

const hasScannedQuery = `SELECT EXISTS(SELECT 1 FROM ` + blocksTableName + ` b1 WHERE b1.block_number = ?) AS found;`

const insertBlockQuery = `
INSERT INTO "` + blocksTableName + `"
	(block_number, block_hash, mined_timestamp, tx_count)
VALUES 
	(?,?,?,?);`

const insertTx = `
INSERT INTO "` + transactionsTableName + `"
	(tx_hash, block_number, tx_from, tx_to, amount, nonce, mined_timestamp, tx_order)
VALUES 
	(?,?,?,?,?,?,?,?);`
