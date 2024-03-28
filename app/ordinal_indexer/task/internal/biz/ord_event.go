package biz

import (
	"context"
	"fmt"
	"github.com/6block/fox_ordinal/pkg/utils"
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

	//IsUsed 如果是transfer-inscribe类型的事件，用于用于标识交易是否已经被使用
	//IsUsed 如果是transfer-transfer类型的事件，用于标识其使用的transfer-inscribe事件的related_ord_transfer_id
	IsUsed     int64 `gorm:"is_used" json:"is_used"`
	EventTime  int64 `gorm:"event_time" json:"event_time"`
	CreateTime int64 `gorm:"create_time" json:"create_time"`
	UpdateTime int64 `gorm:"update_time" json:"update_time"`
}

type IOrdinalEventDAO interface {
	//BatchInsertEvent 批量插入事件
	BatchInsertEvent(ctx context.Context, eventList []OrdEventBO) error

	//GetEventList 查询事件列表
	GetEventList(
		ctx context.Context,
		chainId string,
		collectionId int64,
		inscriptionId string,
		eventStatus string,
		minBlockHeight int64,
		maxBlockHeight int64,
		fromAddressList []string,
		toAddressList []string) ([]OrdEventBO, error)

	//ReorgEvent 重置indexer已同步的事件到newBlockHeight高度,大于该高度的事件以及处理全部删除
	ReorgEvent(
		ctx context.Context,
		chainId string,
		newBlockHeight int64,
	) error

	//UpdateEventStatus 更新事件状态
	UpdateEventStatus(ctx context.Context, eventId int64, collectionId int64, oldStatus, newStatus string) error

	MarkTransferInscribeEventAsUsed(ctx context.Context, relatedOrdTransferId int64) error
}

type OrdinalEventService struct {
	ordinalEventDAO IOrdinalEventDAO
	dataSourceDAO   IDataSourceDAO
	blockDao        IBlockDAO
}

func NewOrdinalEventService(ordinalEventDAO IOrdinalEventDAO, dataSourceDAO IDataSourceDAO, blockDao IBlockDAO) *OrdinalEventService {
	return &OrdinalEventService{
		ordinalEventDAO: ordinalEventDAO,
		dataSourceDAO:   dataSourceDAO,
		blockDao:        blockDao}
}

func (s *OrdinalEventService) SyncEvent(ctx context.Context) (err error) {
	//TODO: 加锁防止并发
	blockInfo, err := s.blockDao.GetBlockInfoByChainId(ctx, s.dataSourceDAO.GetChainId())
	if err != nil {
		return
	}
	syncCount := blockInfo.SyncCountPerTime
	if syncCount < 1 {
		syncCount = 1
	}
	lastBlockWaitForSync := blockInfo.LatestSyncedBlockHeight + syncCount
	if lastBlockWaitForSync > blockInfo.TotalBlockHeight {
		//更新TotalBlockCount
		totalBlockHeight, err := s.dataSourceDAO.GetBlockNumber(ctx)
		if err == nil {
			blockInfo.TotalBlockHeight = totalBlockHeight
		}
	}
	syncCount = utils.Min(syncCount, blockInfo.TotalBlockHeight-blockInfo.LatestSyncedBlockHeight)
	if syncCount < 1 {
		return
	}

	//检查是否存在reorg
	eventList, _ := s.ordinalEventDAO.GetEventList(
		ctx,
		s.dataSourceDAO.GetChainId(),
		0,
		"",
		"",
		blockInfo.LatestSyncedBlockHeight+1,
		0,
		nil,
		nil)

	if len(eventList) > 0 {
		err = s.ordinalEventDAO.ReorgEvent(ctx, s.dataSourceDAO.GetChainId(), blockInfo.LatestSyncedBlockHeight)
		if err != nil {
			fmt.Printf("ReorgEvent failed!, error: %s\n", err.Error())
			return
		}
	}

	list, lastBlockHeight, lastBlockHash, err := s.dataSourceDAO.BatchFetchOrdinalEvent(ctx, blockInfo.LatestSyncedBlockHeight+1, syncCount)
	if err != nil || lastBlockHeight <= blockInfo.LatestSyncedBlockHeight {
		return
	}
	if len(list) > 0 {
		if err = s.ordinalEventDAO.BatchInsertEvent(ctx, list); err != nil {
			return
		}
	}
	blockInfo.LatestSyncedBlockHeight = lastBlockHeight
	blockInfo.BlockHash = lastBlockHash
	err = s.blockDao.UpdateBlockInfo(ctx, blockInfo)
	return err
}

// ListSupportedChainBlock 获取所有支持的链的block同步信息
func (s *OrdinalEventService) ListSupportedChainBlock(ctx context.Context) ([]BlockBO, error) {
	return s.blockDao.ListSupportedChainBlock(ctx)
}

// UpdateEventStatus 更新event状态
func (s *OrdinalEventService) UpdateEventStatus(ctx context.Context, eventId int64, collectionId int64, oldStatus, newStatus string) error {
	return s.ordinalEventDAO.UpdateEventStatus(ctx, eventId, collectionId, oldStatus, newStatus)
}

// GetEventList 从指定链所有未处理事件中获取订阅地址相关事件，保证事件是有序的，保证事件的block高度小于等于最新同步的高度
func (s *OrdinalEventService) GetEventList(
	ctx context.Context,
	chainId string,
	inscriptionId string,
	collectionId int64,
	eventStatus string,
	minBlockHeight int64,
	maxBlockHeight int64,
	fromAddressList []string,
	toAddressList []string) ([]OrdEventBO, error) {
	return s.ordinalEventDAO.GetEventList(
		ctx,
		chainId,
		collectionId,
		inscriptionId,
		eventStatus,
		0,
		maxBlockHeight,
		fromAddressList,
		toAddressList)
}
