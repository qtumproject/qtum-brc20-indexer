package db

type OrdTransfersPO struct {
	ID             int64  `gorm:"AUTO_INCREMENT" json:"id,omitempty"`
	InscriptionId  string `gorm:"inscription_id" json:"inscription_id"`
	BlockHeight    int64  `gorm:"block_height" json:"block_height"`
	OldSatpoint    string `gorm:"old_satpoint" json:"old_satpoint"`
	NewSatpoint    string `gorm:"new_satpoint" json:"new_satpoint"`
	NewPkscript    string `gorm:"new_pkscript" json:"new_pkscript"`
	NewWallet      string `gorm:"new_wallet" json:"new_wallet"`
	SentAsFee      int64  `gorm:"sent_as_fee" json:"sent_as_fee"`
	NewOutputValue int64  `gorm:"new_output_value" json:"new_output_value"`
}

func (OrdTransfersPO) TableName() string {
	return "transfer"
}
