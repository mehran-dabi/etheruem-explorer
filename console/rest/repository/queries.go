package repository

const (
	blocksTableName       = "blocks"
	transactionsTableName = "transactions"
)
const selectLatestBlockQuery = `
SELECT block_number, block_hash, mined_timestamp, tx_count
FROM ` + blocksTableName + `
ORDER BY mined_timestamp DESC
LIMIT 1`

const selectBlockQuery = `
SELECT block_number, block_hash, mined_timestamp, tx_count
FROM ` + blocksTableName + `
WHERE block_number = ?`

const selectAllTxsHashesByBlockIDQuery = `SELECT tx_hash FROM ` + transactionsTableName + ` WHERE block_number = ?`

const selectLatestTxQuery = `
SELECT
	t1.tx_hash,
	t1.block_number,
	t1.tx_from,
	t1.tx_to,
	t1.amount,
	t1.nonce,
	t1.mined_timestamp,
	t1.tx_order
FROM ` + transactionsTableName + ` t1
WHERE
	t1.block_number = (
		SELECT
			b1.block_number
		FROM ` + blocksTableName + ` b1
		ORDER BY
			b1.mined_timestamp DESC
		LIMIT 1)
ORDER BY
	t1.tx_order DESC
LIMIT 1
`

const selectTxQuery = `
SELECT tx_hash, block_number, tx_from, tx_to, amount, nonce, mined_timestamp, tx_order
FROM ` + transactionsTableName + `
WHERE tx_hash = ?`

const selectSumOfAllTxQuery = `
SELECT TOTAL (amount)
FROM ` + transactionsTableName + `
WHERE (block_number BETWEEN ? AND ?);`

const selectAllTxHashQuery = `
SELECT tx_hash
FROM ` + transactionsTableName + `
WHERE (block_number BETWEEN ? AND ?)`
