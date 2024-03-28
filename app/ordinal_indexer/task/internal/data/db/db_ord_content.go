package db

type OrdContentPO struct {
	ID            int64  `gorm:"AUTO_INCREMENT" json:"id,omitempty"`
	InscriptionId string `gorm:"inscription_id" json:"inscription_id"`
	Content       string `gorm:"content" json:"content"`
	TextContent   string `gorm:"text_content" json:"text_content"`
	ContentType   string `gorm:"content_type" json:"content_type"`
	MetaProtocol  string `gorm:"metaprotocol" json:"metaprotocol"`
	BlockHeight   int64  `gorm:"block_height" json:"block_height"`
}

func (OrdContentPO) TableName() string {
	return "ord_content"
}
