package entities

import "time"

// Tx type for blockchain transaction
type Tx struct {
	Hash        string    `json:"hash"`
	BlockNumber int64     `json:"block_number"`
	From        string    `json:"from"`
	To          string    `json:"to"`
	Amount      int64     `json:"amount"`
	Nonce       uint64    `json:"nonce"`
	Timestamp   time.Time `json:"timestamp"`
	Order       int       `json:"order"`
}
