package data

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/task/internal/biz"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/task/internal/config"
	qtum "github.com/6block/fox_ordinal/pkg/qtum_client"
	"time"
)

type QtumDataSourceDAO struct {
	qtumClient *qtum.QtumClient
}

func NewQtumDataSourceDAO(config *config.Config) biz.IDataSourceDAO {
	return &QtumDataSourceDAO{
		qtumClient: qtum.NewQtumClient(
			config.QtumDataSource.Url,
			config.QtumDataSource.AccessToken,
			config.QtumDataSource.NetType),
	}
}

//func (dao *QtumDataSourceDAO) FetchOrdinalEvent(ctx context.Context, blockHeight int64) (eventList []biz.OrdEventBO, err error) {
//	blockInfo, err := dao.qtumClient.GetBlockByHeight(ctx, blockHeight)
//	if err != nil {
//		return eventList, err
//	}
//	return dao.getBRC20OrdinalEventFromBlock(blockInfo)
//}

func (dao *QtumDataSourceDAO) GetBlockNumber(ctx context.Context) (blockHeight int64, err error) {
	blockHeight, err = dao.qtumClient.BlockNumber(ctx)
	return
}

func (dao *QtumDataSourceDAO) BatchFetchOrdinalEvent(
	ctx context.Context,
	startBlockHeight,
	blockCount int64) ([]biz.OrdEventBO, int64, string, error) {
	eventList := make([]biz.OrdEventBO, 0)
	lastBlockHeight := startBlockHeight - 1
	lastBlockHash := ""
	endBlockHeight := startBlockHeight + blockCount - 1
	for height := startBlockHeight; height <= endBlockHeight; height++ {
		blockInfo, err := dao.qtumClient.GetBlockByHeight(ctx, height)
		if err != nil {
			return eventList, lastBlockHeight, lastBlockHash, err
		}
		tempList, err := dao.getBRC20OrdinalEventFromBlock(blockInfo)
		if err != nil {
			return eventList, lastBlockHeight, lastBlockHash, err
		}
		if len(tempList) > 0 {
			eventList = append(eventList, tempList...)
		}
		lastBlockHeight = int64(blockInfo.Height)
		lastBlockHash = blockInfo.Hash
	}
	return eventList, lastBlockHeight, lastBlockHash, nil
}

func (dao *QtumDataSourceDAO) GetChainId() string {
	return "qtum"
}

func (dao *QtumDataSourceDAO) getBRC20OrdinalEventFromBlock(block *qtum.BlockInfo) (eventList []biz.OrdEventBO, err error) {
	eventList = make([]biz.OrdEventBO, 0)
	for _, txInfo := range block.Tx {
		for _, vinInfo := range txInfo.Vin {
			var voutInfo qtum.Vout
			if vinInfo.ScriptSig.Hex == "" && vinInfo.ScriptSig.Asm == "" && vinInfo.Txinwitness != nil {
				if len(vinInfo.Txinwitness) == 3 {
					//从前一个交易中获取fromAddress地址信息
					preTxInfo, err := dao.qtumClient.GetTransactionInfo(context.Background(), vinInfo.Txid, "")
					if err != nil {
						fmt.Println("GetTransactionInfo error, failed to get preTxInfo: " + err.Error())
						continue
					}
					if len(preTxInfo.Vout) <= vinInfo.Vout {
						fmt.Printf("GetTransactionInfo error, preTxInfo not contain index vinInfo.Vout: %d\n ", vinInfo.Vout)
						continue
					}
					bytes, err := hex.DecodeString(vinInfo.Txinwitness[1])
					if err != nil || len(bytes) < 69 {
						continue
					}
					//67 68两个byte以大端形式记录了brc20铭文字符串的长度
					jsonLen := binary.BigEndian.Uint16(bytes[67:69])
					if len(bytes) < int(69+jsonLen) {
						continue
					}
					var content biz.InscriptionBRC20Content
					var vo biz.OrdEventBO
					jsonBytes := bytes[69 : 69+jsonLen]
					err = json.Unmarshal(jsonBytes, &content)
					if err != nil {
						continue
					}
					//TODO: 校验
					//目前只监听brc20协议的铭文事件
					if content.Proto != biz.BRC20_P {
						continue
					}
					fromAddress, err := dao.qtumClient.ExtractQtumAddress(preTxInfo.Vout[vinInfo.Vout])
					if err != nil {
						fmt.Printf("ExtractQtumAddress from %v error: %s\n", preTxInfo.Vout[vinInfo.Vout], err.Error())
					}

					toAddress, err := dao.qtumClient.ExtractQtumAddress(txInfo.Vout[0])
					if err != nil {
						fmt.Printf("ExtractQtumAddress from %v error: %s\n", voutInfo, err.Error())
					}

					vo = biz.OrdEventBO{
						Protocol:        content.Proto,
						ChainId:         dao.GetChainId(),
						EventStatus:     "new",
						Tick:            content.BRC20Tick,
						Amount:          content.BRC20Amount,
						Limit:           content.BRC20Limit,
						Decimal:         content.BRC20Decimal,
						BlockHeight:     int64(block.Height),
						BlockHash:       block.Hash,
						TransactionHash: txInfo.Hash,
						TransactionId:   txInfo.Txid,
						Max:             content.BRC20Max,
						SourceAddress:   fromAddress,
						TargetAddress:   toAddress,
						CreateTime:      time.Now().Unix(),
						UpdateTime:      time.Now().Unix(),
					}
					if len(vo.Decimal) == 0 {
						vo.Decimal = biz.DEFAULT_DECIMAL_18
					}
					eventList = append(eventList, vo)
				}
			}
		}
	}
	return
}
