package db

type OrdNumberToIdPO struct {
	ID               int64   `gorm:"AUTO_INCREMENT" json:"id,omitempty"`
	HolderAddress    string  `gorm:"holder_address" json:"holderAddress"`
	ChainId          string  `gorm:"chain_id" json:"chainId"`
	CollectionId     int64   `gorm:"collection_id" json:"collectionId"`
	EventId          int64   `gorm:"event_id" json:"eventId"`
	TransferType     int64   `gorm:"transfer_type" json:"transfer_type"`     //0 支出； 1 收入
	TransferStatus   int64   `gorm:"transfer_status" json:"transfer_status"` //0 switch off; 1 switch on
	IsConfirmDeposit string  `gorm:"is_confirm_deposit" json:"isConfirmDeposit"`
	Amount           float64 `gorm:"amount" json:"amount"`
	TotalAmount      float64 `gorm:"total_amount" json:"totalAmount"`
	Remark           string  `gorm:"remark" json:"remark"`
	BlockHeight      int64   `gorm:"block_height" json:"block_height"`
	TransactionHash  string  `gorm:"transaction_hash" json:"transaction_hash"`
	TransactionId    string  `gorm:"transaction_id" json:"transaction_id"`
	CreateTime       int64   `gorm:"create_time" json:"createTime"`
	UpdateTime       int64   `gorm:"update_time" json:"updateTime"`
}

func (OrdNumberToIdPO) TableName() string {
	return "ord_number_to_id"
}
