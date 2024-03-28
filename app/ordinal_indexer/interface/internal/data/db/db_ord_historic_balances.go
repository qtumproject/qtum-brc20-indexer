package db

type OrdHistoricBalancesPO struct {
	ID               int64   `gorm:"AUTO_INCREMENT" json:"id,omitempty"`
	ChainId          string  `gorm:"chain_id" json:"chainId"`
	CollectionId     int64   `gorm:"collection_id" json:"collectionId"`
	Tick             string  `gorm:"tick" json:"tick"`
	Pkscript         string  `gorm:"pkscript" json:"pkscript"`
	WalletAddress    string  `gorm:"wallet_address" json:"wallet_address"`
	BlockHeight      int64   `gorm:"block_height" json:"block_height"`
	TransactionHash  string  `gorm:"transaction_hash" json:"transaction_hash"`
	TransactionId    string  `gorm:"transaction_id" json:"transaction_id"`
	EventId          int64   `gorm:"event_id" json:"eventId"`
	ChangeType       int64   `gorm:"change_type" json:"change_type"` //0 支出； 1 收入
	OverallAmount    float64 `gorm:"overall_amount" json:"overall_amount"`
	AvailableAmount  float64 `gorm:"available_amount" json:"available_amount"`
	OverallBalance   float64 `gorm:"overall_balance" json:"overall_balance"`
	AvailableBalance float64 `gorm:"available_balance" json:"available_balance"`
	CreateTime       int64   `gorm:"create_time" json:"createTime"`
	UpdateTime       int64   `gorm:"update_time" json:"updateTime"`
}

func (OrdHistoricBalancesPO) TableName() string {
	return "ord_historic_balances"
}
