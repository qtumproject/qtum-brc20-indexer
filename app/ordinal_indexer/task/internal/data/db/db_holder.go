package db

type HolderPO struct {
	ID                int64   `gorm:"AUTO_INCREMENT" json:"id,omitempty"`
	Pkscript          string  `gorm:"pkscript" json:"pkscript"`
	ChainId           string  `gorm:"chain_id" json:"chain_id"`
	CollectionId      int64   `gorm:"collection_id" json:"collection_id"`
	Tick              string  `gorm:"tick" json:"tick"`
	Address           string  `gorm:"address" json:"address"`
	LastBlockHeight   int64   `gorm:"last_block_height" json:"last_block_height"`
	LastTransactionId string  `gorm:"last_transaction_id" json:"last_transaction_id"`
	OverallBalance    float64 `gorm:"overall_balance" json:"overall_balance"`
	AvailableBalance  float64 `gorm:"available_balance" json:"available_balance"`
	CreateTime        int64   `gorm:"create_time" json:"create_time"`
	UpdateTime        int64   `gorm:"update_time" json:"update_time"`
}

func (HolderPO) TableName() string {
	return "holder"
}
