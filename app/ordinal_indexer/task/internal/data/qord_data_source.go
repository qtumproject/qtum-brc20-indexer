package data

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/task/internal/biz"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/task/internal/config"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/task/internal/data/db"
	qtum "github.com/6block/fox_ordinal/pkg/qtum_client"
	"strings"
	"time"
	"unicode/utf8"
)

type QordDataSourceDAO struct {
	data                   *Data
	firstInscriptionHeight int64
	qtumClient             *qtum.QtumClient
}

func NewQordDataSourceDAO(data *Data, config *config.Config) biz.IDataSourceDAO {
	return &QordDataSourceDAO{
		data:                   data,
		firstInscriptionHeight: 3547171, //TODO 走配置
		qtumClient: qtum.NewQtumClient(
			config.QtumDataSource.Url,
			config.QtumDataSource.AccessToken,
			config.QtumDataSource.NetType),
	}
}

func (dao *QordDataSourceDAO) GetBlockNumber(ctx context.Context) (blockHeight int64, err error) {
	//blockHeight, err = dao.qtumClient.BlockNumber(ctx)
	list := make([]db.OrdBlockHashesPO, 0)
	err = dao.data.db.SqlDB.Model(&db.OrdBlockHashesPO{}).Order("block_height desc").Limit(1).Find(&list).Error
	if err != nil {
		return
	}
	if len(list) == 0 {
		blockHeight = dao.firstInscriptionHeight
	} else {
		blockHeight = list[0].BlockHeight
	}
	return
}

func (dao *QordDataSourceDAO) GetBlockHashByHeight(ctx context.Context, blockHeight int64) (blockHash string, err error) {
	var blockInfo db.OrdBlockHashesPO
	err = dao.data.db.SqlDB.Model(&db.OrdBlockHashesPO{}).Where("block_height = ?", blockHeight).First(&blockInfo).Error
	if err != nil {
		return
	}
	blockHash = blockInfo.BlockHash
	return
}

func (dao *QordDataSourceDAO) ListBlockInfo(ctx context.Context, minBlockHeight, count int64) ([]db.OrdBlockHashesPO, error) {
	var blockInfoList []db.OrdBlockHashesPO
	err := dao.data.db.SqlDB.Model(&db.OrdBlockHashesPO{}).Where("block_height >= ?", minBlockHeight).
		Order("block_height asc").Limit(int(count)).
		Find(&blockInfoList).Error
	return blockInfoList, err
}

func (dao *QordDataSourceDAO) BatchFetchOrdinalEvent(
	ctx context.Context,
	startBlockHeight,
	blockCount int64) ([]biz.OrdEventBO, int64, string, error) {

	eventList := make([]biz.OrdEventBO, 0)
	lastBlockHeight := startBlockHeight - 1
	lastBlockHash := ""
	endBlockHeight := startBlockHeight + blockCount - 1
	list, err := dao.ListBlockInfo(ctx, startBlockHeight, blockCount+1)
	heightToHashMap := make(map[int64]string)
	if err != nil {
		return eventList, lastBlockHeight, lastBlockHash, err
	}
	for _, blockInfo := range list {
		heightToHashMap[blockInfo.BlockHeight] = blockInfo.BlockHash
	}
	for height := startBlockHeight; height <= endBlockHeight; height++ {

		currentBlockHash, ok := heightToHashMap[height]
		if !ok {
			return eventList, lastBlockHeight, lastBlockHash, err
		}
		tempList, err := dao.GetBRC20OrdinalEvent(ctx, height)
		if err != nil {
			return eventList, lastBlockHeight, lastBlockHash, err
		}

		blockInfo, err := dao.qtumClient.GetBlockVerboseByHeight(ctx, height)

		for index, _ := range tempList {
			tempList[index].BlockHash = currentBlockHash
			tempList[index].EventTime = blockInfo.Time
		}
		if len(tempList) > 0 {
			eventList = append(eventList, tempList...)
		}

		lastBlockHeight = height
		lastBlockHash = currentBlockHash
		if height%100 == 0 {
			fmt.Printf("【BatchFetchOrdinalEvent】 block %d event collected, current event count: %d\n", height, len(eventList))
		}

	}
	return eventList, lastBlockHeight, lastBlockHash, nil
}

func (dao *QordDataSourceDAO) GetChainId() string {
	return "qtum"
}

func (dao *QordDataSourceDAO) GetValidOrdEventList(
	ctx context.Context,
	chainId string,
	protocol string,
	tick string,
	inscriptionId string) (eventList []biz.OrdEventBO, err error) {

	query := dao.data.db.SqlDB.Model(&db.OrdEventPO{}).
		Where("chain_id = ?", chainId).
		Where("protocol = ?", protocol).
		Where("tick = ?", tick).
		Where("inscription_id = ?", inscriptionId).
		Where("event_status != ?", "invalid").
		Order("id ASC")
	err = query.Find(&eventList).Error
	return
}

func (dao *QordDataSourceDAO) GetBRC20OrdinalEvent(ctx context.Context, height int64) (eventList []biz.OrdEventBO, err error) {
	eventList = make([]biz.OrdEventBO, 0)
	type QordTransferEventInfo struct {
		ID             int64  `gorm:"AUTO_INCREMENT" json:"id,omitempty"`
		InscriptionId  string `gorm:"inscription_id" json:"inscription_id"`
		BlockHeight    int64  `gorm:"block_height" json:"block_height"`
		OldSatpoint    string `gorm:"old_satpoint" json:"old_satpoint"`
		NewSatpoint    string `gorm:"new_satpoint" json:"new_satpoint"`
		NewPkscript    string `gorm:"new_pkscript" json:"new_pkscript"`
		NewWallet      string `gorm:"new_wallet" json:"new_wallet"`
		SentAsFee      int64  `gorm:"sent_as_fee" json:"sent_as_fee"`
		NewOutputValue int64  `gorm:"new_output_value" json:"new_output_value"`
		Content        string `gorm:"content" json:"content"`
		ContentType    string `gorm:"content_type" json:"content_type"`
	}
	qordEventInfoList := make([]QordTransferEventInfo, 0)
	query :=
		`SELECT ot.id, ot.inscription_id, ot.old_satpoint, ot.new_satpoint, ot.new_pkscript, ot.new_wallet, ot.sent_as_fee, oc.content, oc.content_type
			FROM ord_transfers ot
			LEFT JOIN ord_content oc ON ot.inscription_id = oc.inscription_id
			LEFT JOIN ord_number_to_id onti ON ot.inscription_id = onti.inscription_id
			WHERE ot.block_height = {{BlockHeight}}
			AND onti.cursed_for_brc20 = false
			AND oc.content is not null AND oc.content like '{%'
			ORDER BY ot.id asc;`
	query = strings.Replace(query, "{{BlockHeight}}", fmt.Sprintf("%d", height), -1)
	if err = dao.data.db.SqlDB.Raw(query).Find(&qordEventInfoList).Error; err != nil {
		return
	}

	for _, qordEventInfo := range qordEventInfoList {
		if qordEventInfo.SentAsFee == 1 && qordEventInfo.OldSatpoint == "" {
			//inscribed as fee
			continue
		}

		if qordEventInfo.ContentType == "" {
			//invalid inscription
			continue
		}
		//decode content_type filed
		bytes, _ := hex.DecodeString(qordEventInfo.ContentType)
		decodedContentType := string(bytes)
		contentType := strings.Split(decodedContentType, ";")[0]
		if contentType != "application/json" && contentType != "text/plain" {
			//invalid inscription
			fmt.Printf("【GetBRC20OrdinalEvent】invalid event %v, invalid contentType %s\n", qordEventInfo, contentType)
			continue
		}
		//qordEventInfo.Content
		var content biz.InscriptionBRC20Content
		err = json.Unmarshal([]byte(qordEventInfo.Content), &content)
		if err != nil {
			fmt.Printf("【GetBRC20OrdinalEvent】invalid event %v, parse content error: %s\n", qordEventInfo, err.Error())
			continue
		}

		if content.Proto != "brc-20" || content.BRC20Tick == "" || content.Operation == "" {
			fmt.Printf("【GetBRC20OrdinalEvent】invalid event %v, invalid content : %v\n", qordEventInfo, content)
			continue
		}
		content.BRC20Tick = strings.ToLower(content.BRC20Tick)
		if utf8.RuneCountInString(content.BRC20Tick) != 4 {
			fmt.Printf("【GetBRC20OrdinalEvent】invalid event %v, invalid tick len: %s\n", qordEventInfo, content.BRC20Tick)
			continue
		}

		vo := biz.OrdEventBO{
			Protocol:             content.Proto,
			ChainId:              dao.GetChainId(),
			BlockHeight:          height,
			TransactionId:        strings.Split(qordEventInfo.NewSatpoint, ":")[0],
			EventStatus:          "new",
			Tick:                 content.BRC20Tick,
			InscriptionId:        qordEventInfo.InscriptionId,
			RelatedOrdTransferId: qordEventInfo.ID,
			Amount:               content.BRC20Amount,
			Limit:                content.BRC20Limit,
			Decimal:              content.BRC20Decimal,
			Max:                  content.BRC20Max,
			CreateTime:           time.Now().Unix(),
			UpdateTime:           time.Now().Unix(),
		}
		if len(vo.Decimal) == 0 {
			vo.Decimal = biz.DEFAULT_DECIMAL_18
		}

		//处理 event_type字段
		if content.Operation == "deploy" && qordEventInfo.OldSatpoint == "" {
			vo.EventType = biz.DeployInscribe
			vo.TargetPkscript = qordEventInfo.NewPkscript
			vo.TargetAddress = qordEventInfo.NewWallet
		} else if content.Operation == "mint" && qordEventInfo.OldSatpoint == "" {
			vo.EventType = biz.MintInscribe
			vo.TargetPkscript = qordEventInfo.NewPkscript
			vo.TargetAddress = qordEventInfo.NewWallet

		} else if content.Operation == "transfer" {
			if qordEventInfo.OldSatpoint == "" {
				vo.EventType = biz.TransferInscribe
				vo.SourcePkscript = qordEventInfo.NewPkscript
				vo.SourceAddress = qordEventInfo.NewWallet
			} else {
				vo.EventType = biz.TransferTransfer

				if qordEventInfo.SentAsFee == 1 {
					//transfer_transfer_spend_to_fee
					vo.TargetPkscript = ""
					vo.TargetAddress = ""
				} else {
					//transfer_transfer_normal
					vo.TargetPkscript = qordEventInfo.NewPkscript
					vo.TargetAddress = qordEventInfo.NewWallet
				}

				ordEventList, err := dao.GetValidOrdEventList(
					ctx,
					dao.GetChainId(),
					"brc-20",
					content.BRC20Tick,
					qordEventInfo.InscriptionId,
				)

				//处理TransferTransfer引用的TransferInscribe交易在当前区块中的场景
				for _, event := range eventList {
					if event.Tick == content.BRC20Tick && event.InscriptionId == qordEventInfo.InscriptionId {
						ordEventList = append(ordEventList, event)
					}
				}

				if err != nil || len(ordEventList) == 0 {
					if len(ordEventList) == 0 {
						fmt.Printf("【GetBRC20OrdinalEvent】invalid Transfer-Transfer event %v, no InscriptionId found\n", qordEventInfo)
					} else {
						fmt.Printf("【GetBRC20OrdinalEvent】invalid Transfer-Transfer event %v, GetValidOrdEventList error: %s\n", qordEventInfo, err.Error())
					}

					continue
				}
				//check if InscriptionId is used or invalid
				isValid := false

				for _, event := range ordEventList {
					if event.EventType == biz.TransferTransfer {
						//InscriptionId used, invalid
						fmt.Printf("【GetBRC20OrdinalEvent】invalid Transfer-Transfer event %v, InscriptionId used\n", qordEventInfo)
						break
					} else if event.EventType == biz.TransferInscribe {
						//first TransferInscribe event is valid
						isValid = true
						vo.IsUsed = event.RelatedOrdTransferId
						vo.SourcePkscript = event.SourcePkscript
						vo.SourceAddress = event.SourceAddress
						break
					}
				}
				if !isValid {
					fmt.Printf("【GetBRC20OrdinalEvent】invalid Transfer-Transfer event %v, InscriptionId not found\n", qordEventInfo)
					continue
				}
			}
		} else {
			continue
		}
		eventList = append(eventList, vo)
	}
	return
}
