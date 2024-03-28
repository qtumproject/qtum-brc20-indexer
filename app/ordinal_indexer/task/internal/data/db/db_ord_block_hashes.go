package db

type OrdBlockHashesPO struct {
	ID          int64  `gorm:"AUTO_INCREMENT" json:"id,omitempty"`
	BlockHash   string `gorm:"block_hash" json:"block_hash"`
	BlockHeight int64  `gorm:"block_height" json:"block_height"` //总区块数量
}

func (OrdBlockHashesPO) TableName() string {
	return "block_hashes"
}
