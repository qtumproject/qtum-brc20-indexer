package biz

import (
	"context"
	"fmt"
)

type CollectionBO struct {
	ID              int64   `gorm:"AUTO_INCREMENT" json:"id,omitempty"`
	ChainId         string  `gorm:"chain_id" json:"chainId"`
	EventId         int64   `gorm:"event_id" json:"event_id"`
	Protocol        string  `gorm:"protocol" json:"protocol"`
	Tick            string  `gorm:"tick" json:"tick"`
	MaxAmount       float64 `gorm:"max_amount" json:"max_amount"`               //最大供应，当 op = deploy 的时候用到
	MintLimitAmount float64 `gorm:"mint_limit_amount" json:"mint_limit_amount"` //单次铸造时的数量上限，当 op = mint 的时候用到
	TotalMinted     float64 `gorm:"total_minted" json:"total_minted"`           //已铸造数量
	MintTimes       int64   `gorm:"mint_times" json:"mint_times"`               //已铸造次数
	Decimal         int64   `gorm:"decimal" json:"decimal"`

	BlockHeight       int64  `gorm:"block_height" json:"block_height"`
	TransactionId     string `gorm:"transaction_id" json:"transaction_id"`
	LastBlockHeight   int64  `gorm:"last_block_height" json:"last_block_height"`
	LastTransactionId string `gorm:"last_transaction_id" json:"last_transaction_id"`
	TransactionHash   string `gorm:"transaction_hash" json:"transaction_hash"`
	DeployTime        int64  `gorm:"deploy_time" json:"deploy_time"`
	CreateTime        int64  `gorm:"create_time" json:"create_time"`
	UpdateTime        int64  `gorm:"update_time" json:"update_time"`
}

type InscriptionBRC20Content struct {
	Proto        string `json:"p,omitempty"`
	Operation    string `json:"op,omitempty"`
	BRC20Tick    string `json:"tick,omitempty"`
	BRC20Max     string `json:"max,omitempty"`
	BRC20Amount  string `json:"amt,omitempty"`
	BRC20Limit   string `json:"lim,omitempty"` // option
	BRC20Decimal string `json:"dec,omitempty"` // option
}

type ICollectionDAO interface {
	GetCollectionInfo(ctx context.Context, chainId, protocol, tick string) (CollectionBO, error)
	CreateCollection(ctx context.Context, collectionInfo *CollectionBO) error
	UpdateCollectionInfo(ctx context.Context, id int64, bo CollectionBO) error
}

type CollectionService struct {
	collectionDAO ICollectionDAO
}

func NewCollectionService(repo ICollectionDAO) *CollectionService {
	return &CollectionService{collectionDAO: repo}
}

func (cc *CollectionService) GetCollection(ctx context.Context, chainId string, protocol string, tick string) (CollectionBO, error) {
	return cc.collectionDAO.GetCollectionInfo(ctx, chainId, protocol, tick)
}

func (cc *CollectionService) CreateCollection(ctx context.Context, bo *CollectionBO) error {
	err := cc.collectionDAO.CreateCollection(ctx, bo)
	return err
}

func (cc *CollectionService) UpdateCollectionInfo(ctx context.Context, bo CollectionBO) error {
	if bo.ID <= 0 {
		return fmt.Errorf("UpdateCollectionInfo failed, id can not be empty")
	}
	err := cc.collectionDAO.UpdateCollectionInfo(ctx, bo.ID, bo)
	return err
}
