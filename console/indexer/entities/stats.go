package entities

// Stats - type to return the transaction statistics
type Stats struct {
	Txs         []string `json:"txs"`
	TotalAmount float64  `json:"total_amount"`
}
