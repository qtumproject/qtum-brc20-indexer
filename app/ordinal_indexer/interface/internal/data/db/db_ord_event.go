package db

type OrdEventPO struct {
	ID                   int64  `gorm:"AUTO_INCREMENT" json:"id,omitempty"`
	Protocol             string `gorm:"protocol" json:"protocol"`
	ChainId              string `gorm:"chain_id" json:"chain_id"`
	Tick                 string `gorm:"tick" json:"tick"`
	CollectionId         int64  `gorm:"collection_id" json:"collection_id"`
	EventType            int64  `gorm:"event_type" json:"event_type"`
	BlockHeight          int64  `gorm:"block_height" json:"block_height"`
	BlockHash            string `gorm:"block_hash" json:"block_hash"`
	TransactionHash      string `gorm:"transaction_hash" json:"transaction_hash"`
	TransactionId        string `gorm:"transaction_id" json:"transaction_id"`
	InscriptionId        string `gorm:"inscription_id" json:"inscription_id"`
	RelatedOrdTransferId string `gorm:"related_ord_transfer_id" json:"related_ord_transfer_id"`
	EventStatus          string `gorm:"event_status" json:"event_status"`
	SourceAddress        string `gorm:"source_address" json:"source_address"`
	SourcePkscript       string `gorm:"source_pkscript" json:"source_pkscript"`
	TargetAddress        string `gorm:"target_address" json:"target_address"`
	TargetPkscript       string `gorm:"target_pkscript" json:"target_pkscript"`
	Max                  string `gorm:"max" json:"max"`
	Limit                string `gorm:"limit" json:"limit"`
	Amount               string `gorm:"amount" json:"amount"`
	Decimal              string `gorm:"decimal" json:"decimal"`
	CallData             string `gorm:"call_data" json:"call_data"`
	IsUsed               int64  `gorm:"is_used" json:"is_used"` //用于标识transfer-inscribe类型的交易是否已经被使用
	EventTime            int64  `gorm:"event_time" json:"event_time"`

	CreateTime int64 `gorm:"create_time" json:"create_time"`
	UpdateTime int64 `gorm:"update_time" json:"update_time"`
}

func (OrdEventPO) TableName() string {
	return "ord_event"
}
