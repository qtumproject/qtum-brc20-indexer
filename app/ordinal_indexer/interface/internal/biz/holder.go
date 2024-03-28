package biz

import (
	"context"
)

type IHolderDAO interface {
	GetHolderInfo(ctx context.Context, address string, chainId string, collectionId int64) (HolderBO, error)
	GetHolderListByCollectionId(ctx context.Context, chainId string, collectionId int64) ([]HolderBO, error)
	GetHolderCollectionList(ctx context.Context, address, chainId string) ([]HolderBO, error)
	GetHolderListByChainId(ctx context.Context, chainId string) ([]HolderBO, error)
	UpdateAmount(ctx context.Context, id int64, amount float64) error
}

type HolderService struct {
	repo IHolderDAO
}

type HolderBO struct {
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

func NewHolderService(repo IHolderDAO) *HolderService {
	return &HolderService{repo: repo}
}

func (uc *HolderService) GetAmount(ctx context.Context, address string, chainId string, collectionId int64) (balance float64, err error) {
	holderInfo, err := uc.repo.GetHolderInfo(ctx, address, chainId, collectionId)
	if err != nil {
		return
	}
	return holderInfo.OverallBalance, nil
}

func (uc *HolderService) UpdateAmount(ctx context.Context, id int64, amount float64) error {
	err := uc.repo.UpdateAmount(ctx, id, amount)
	return err
}

func (uc *HolderService) GetHolderListByChainId(ctx context.Context, chainId string) (list []HolderBO, err error) {
	list, err = uc.repo.GetHolderListByChainId(ctx, chainId)
	return
}

func (uc *HolderService) GetHolderListByCollectionId(ctx context.Context, chainId string, collectionId int64) ([]HolderBO, error) {
	list, err := uc.repo.GetHolderListByCollectionId(ctx, chainId, collectionId)
	if err != nil {
		return list, err
	}
	realList := make([]HolderBO, 0)
	for _, record := range list {
		if record.AvailableBalance == 0 && record.OverallBalance == 0 {
			continue
		}
		realList = append(realList, record)
	}
	return realList, err
}

func (uc *HolderService) GetHolderCollectionList(ctx context.Context, address, chainId string) (bo []HolderBO, err error) {
	return uc.repo.GetHolderCollectionList(ctx, address, chainId)
}
