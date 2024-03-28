package biz

import (
	"context"
)

type OrdinalEventType int

const (
	DeployInscribe OrdinalEventType = iota
	MintInscribe
	TransferInscribe
	TransferTransfer
)

const (
	Brc20Protocol = "brc-20"
)

var OrdinalEventTypeMap = map[OrdinalEventType]string{
	DeployInscribe:   "deploy-inscribe",
	MintInscribe:     "mint-inscribe",
	TransferInscribe: "transfer-inscribe",
	TransferTransfer: "transfer-transfer",
}

var OrdinalEventTypeReverseMap = map[string]OrdinalEventType{
	"deploy-inscribe":   DeployInscribe,
	"mint-inscribe":     MintInscribe,
	"transfer-inscribe": TransferInscribe,
	"transfer-transfer": TransferTransfer,
}

type OrdEventBO struct {
	ID                   int64            `gorm:"AUTO_INCREMENT" json:"id,omitempty"`
	Protocol             string           `gorm:"protocol" json:"protocol"`
	ChainId              string           `gorm:"chain_id" json:"chain_id"`
	Tick                 string           `gorm:"tick" json:"tick"`
	CollectionId         int64            `gorm:"collection_id" json:"collection_id"`
	EventType            OrdinalEventType `gorm:"event_type" json:"event_type"`
	BlockHeight          int64            `gorm:"block_height" json:"block_height"`
	BlockHash            string           `gorm:"block_hash" json:"block_hash"`
	TransactionHash      string           `gorm:"transaction_hash" json:"transaction_hash"`
	TransactionId        string           `gorm:"transaction_id" json:"transaction_id"`
	InscriptionId        string           `gorm:"inscription_id" json:"inscription_id"`
	RelatedOrdTransferId int64            `gorm:"related_ord_transfer_id" json:"related_ord_transfer_id"`
	EventStatus          string           `gorm:"event_status" json:"event_status"`
	SourceAddress        string           `gorm:"source_address" json:"source_address"`
	SourcePkscript       string           `gorm:"source_pkscript" json:"source_pkscript"`
	TargetAddress        string           `gorm:"target_address" json:"target_address"`
	TargetPkscript       string           `gorm:"target_pkscript" json:"target_pkscript"`
	Max                  string           `gorm:"max" json:"max"`
	Limit                string           `gorm:"limit" json:"limit"`
	Amount               string           `gorm:"amount" json:"amount"`
	Decimal              string           `gorm:"decimal" json:"decimal"`
	CallData             string           `gorm:"call_data" json:"call_data"`
	IsUsed               int64            `gorm:"is_used" json:"is_used"` //用于标识transfer-inscribe类型的交易是否已经被使用
	EventTime            int64            `gorm:"event_time" json:"event_time"`
	CreateTime           int64            `gorm:"create_time" json:"create_time"`
	UpdateTime           int64            `gorm:"update_time" json:"update_time"`
}

type IOrdinalEventDAO interface {

	//GetEventList 查询事件列表
	GetEventList(
		ctx context.Context,
		chainId string,
		collectionId int64,
		inscriptionId string,
		eventStatus string,
		maxBlockHeight int64,
		fromAddressList []string,
		toAddressList []string) ([]OrdEventBO, error)
}

type OrdinalEventService struct {
	ordinalEventDAO IOrdinalEventDAO
	blockDao        IBlockDAO
}

func NewOrdinalEventService(ordinalEventDAO IOrdinalEventDAO, blockDao IBlockDAO) *OrdinalEventService {
	return &OrdinalEventService{
		ordinalEventDAO: ordinalEventDAO,
		blockDao:        blockDao}
}

// ListSupportedChainBlock 获取所有支持的链的block同步信息
func (s *OrdinalEventService) ListSupportedChainBlock(ctx context.Context) ([]BlockBO, error) {
	return s.blockDao.ListSupportedChainBlock(ctx)
}

// FetchEventList 从指定链所有未处理事件中获取订阅地址相关事件，保证事件是有序的，保证事件的block高度小于等于最新同步的高度
func (s *OrdinalEventService) FetchEventList(
	ctx context.Context,
	chainId string,
	inscriptionId string,
	collectionId int64,
	eventStatus string,
	maxBlockHeight int64,
	fromAddressList []string,
	toAddressList []string) ([]OrdEventBO, error) {
	return s.ordinalEventDAO.GetEventList(ctx, chainId, collectionId, inscriptionId, eventStatus, maxBlockHeight, fromAddressList, toAddressList)
}
