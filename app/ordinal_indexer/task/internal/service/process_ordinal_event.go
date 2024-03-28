package service

import (
	"context"
	"fmt"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/task/internal/biz"
	"github.com/6block/fox_ordinal/pkg/decimal"
	"math"
	"strconv"
	"strings"
	"time"
)

//TODO:场景2，新增订阅地址数据处理：从所有事件中获取订阅地址相关事件，依次处理并（落库transfer表&更新用户余额信息）

// ProcessOrdinalEventJob 场景1，常规同步block时的数据处理：从所有未处理状态的事件中获取订阅地址相关的事件，依次处理并（落库transfer表&更新用户余额信息）
func (svc OrdinalIndexerTaskServer) ProcessOrdinalEventJob() {
	ctx := context.Background()
	//1.1 获取系统当前支持的所有链
	chainBlockList, err := svc.ordinalEventService.ListSupportedChainBlock(ctx)
	if err != nil {
		fmt.Println("ProcessOrdinalEvent failed, ListSupportedChainBlock report error: " + err.Error())
		return
	}
	//1.2 逐个链查询数据库中的事件并处理：
	for _, chainBlock := range chainBlockList {
		//1.2.0 上分布式锁，防止并发更新, 获取锁失败则直接continue跳到下一个链
		//1.2.1 获取当前链的订阅地址
		//TODO: 添加缓存
		//holderList, err := svc.holderService.GetHolderListByChainId(ctx, chainBlock.ChainId)
		//if err != nil {
		//	fmt.Println("ProcessOrdinalEvent failed, GetHolderListByChainId report error: " + err.Error())
		//	continue
		//}
		//var subscribedAddress map[string]bool
		//for _, holder := range holderList {
		//	subscribedAddress[holder.Address] = true
		//}
		//1.2.2 从当前链所有未处理事件中获取未处理相关事件(保证事件是有序的，保证事件的block高度小于等于最新同步区块高度)
		eventList, err := svc.ordinalEventService.GetEventList(
			ctx,
			chainBlock.ChainId,
			"",
			0,
			"new",
			0,
			chainBlock.LatestSyncedBlockHeight,
			nil,
			nil)
		if err != nil {
			fmt.Println("ProcessOrdinalEvent failed, GetUnhandledEventList report error: " + err.Error())
			continue
		}

		//1.2.3 依次处理事件并更新未处理事件状态为已处理
		for _, event := range eventList {
			err = svc.ProcessOrdinalEvent(ctx, event)
			if err != nil {
				fmt.Printf("ProcessOrdinalEvent %v failed, error: %s\n", event, err.Error())
			}
		}
		//1.2.4 释放分布式锁
	}

	return
}
func (svc OrdinalIndexerTaskServer) ProcessOrdinalEvent(
	ctx context.Context, event biz.OrdEventBO) error {

	switch event.EventType {
	case biz.DeployInscribe:
		return svc.processDeployEvent(ctx, event)
	case biz.MintInscribe:
		return svc.processMintEvent(ctx, event)
	case biz.TransferInscribe, biz.TransferTransfer:
		return svc.processTransferEvent(ctx, event)
	default:
		return fmt.Errorf("ProcessOrdinalEvent failed, event op %s not define", event.EventType)
	}
}

func (svc OrdinalIndexerTaskServer) processDeployEvent(ctx context.Context, event biz.OrdEventBO) error {
	switch event.Protocol {
	case biz.Brc20Protocol:
		return svc.DeployBRC20OrdinalCollection(ctx, event)
	default:
		return fmt.Errorf("processDeployEvent failed, event Protocol %s not define", event.Protocol)
	}
	return nil
}

func (svc OrdinalIndexerTaskServer) processMintEvent(
	ctx context.Context,
	event biz.OrdEventBO,
) error {
	switch event.Protocol {
	case biz.Brc20Protocol:
		return svc.MintBRC20OrdinalCollection(ctx, event)
	default:
		return fmt.Errorf("ProcessOrdinalEvent failed, event Protocol %s not define", event.Protocol)
	}
	return nil
}

func (svc OrdinalIndexerTaskServer) processTransferEvent(
	ctx context.Context,
	event biz.OrdEventBO,
) error {
	switch event.Protocol {
	case biz.Brc20Protocol:
		return svc.TransferBRC20OrdinalCollection(ctx, event)
	default:
		return fmt.Errorf("ProcessOrdinalEvent failed, event Protocol %s not define", event.Protocol)
	}

	return nil
}

func (svc OrdinalIndexerTaskServer) MintBRC20OrdinalCollection(
	ctx context.Context, eventInfo biz.OrdEventBO) error {

	bo := biz.OrdHistoricBalancesBO{
		ChainId:         eventInfo.ChainId,
		EventId:         eventInfo.ID,
		ChangeType:      1,
		TransactionHash: eventInfo.TransactionHash,
		TransactionId:   eventInfo.TransactionId,
		BlockHeight:     eventInfo.BlockHeight,
		CreateTime:      time.Now().Unix(),
		UpdateTime:      time.Now().Unix(),
	}
	var transferAmountValue float64

	// check tick
	uniqueLowerTicker := strings.ToLower(eventInfo.Tick)
	collection, err := svc.collectionService.GetCollection(ctx, eventInfo.ChainId, eventInfo.Protocol, uniqueLowerTicker)
	if err != nil {
		return err
	}
	if collection.Tick != uniqueLowerTicker {
		//ticker not exist
		fmt.Printf("MintBRC20OrdinalCollection mint event %v, but ticker not exist\n", eventInfo)
		svc.ordinalEventService.UpdateEventStatus(ctx, eventInfo.ID, eventInfo.CollectionId, eventInfo.EventStatus, "invalid")
		return nil
	}
	if collection.TotalMinted >= collection.MaxAmount {
		//mint ended
		fmt.Printf("MintBRC20OrdinalCollection event %v, but collection mint ended", eventInfo)
		return svc.ordinalEventService.UpdateEventStatus(ctx, eventInfo.ID, eventInfo.CollectionId, eventInfo.EventStatus, "invalid")
	}
	bo.CollectionId = collection.ID
	bo.Tick = uniqueLowerTicker
	eventInfo.Tick = uniqueLowerTicker
	eventInfo.CollectionId = collection.ID

	// check mint amount
	if amount, precision, parseErr := decimal.NewDecimalFromString(eventInfo.Amount); parseErr != nil {
		// amount invalid
		fmt.Printf("MintBRC20OrdinalCollection event %v, but amount parse err. amount: '%s', error: %s\n",
			eventInfo,
			eventInfo.Amount,
			parseErr.Error(),
		)
		svc.ordinalEventService.UpdateEventStatus(ctx, eventInfo.ID, eventInfo.CollectionId, eventInfo.EventStatus, "invalid")
		return nil
	} else {
		if amount.Sign() <= 0 || precision > int(collection.Decimal) {
			fmt.Printf("MintBRC20OrdinalCollection event %v, but amount invalid.\n",
				eventInfo)
			svc.ordinalEventService.UpdateEventStatus(ctx, eventInfo.ID, eventInfo.CollectionId, eventInfo.EventStatus, "invalid")
			return nil
		}
		if amount.Float64() > collection.MintLimitAmount {
			//mint too much
			fmt.Printf("MintBRC20OrdinalCollection event %v, but mint too much. amount %s out of limit %f.\n",
				eventInfo, amount.String(), collection.MintLimitAmount)
			svc.ordinalEventService.UpdateEventStatus(ctx, eventInfo.ID, eventInfo.CollectionId, eventInfo.EventStatus, "invalid")
			return nil
		}
		transferAmountValue = amount.Float64()
	}
	//1. 更新token余额并判断实际mint成功的数量
	newTotalMinted := math.Min(collection.MaxAmount, transferAmountValue+collection.TotalMinted)
	realMintedValue := newTotalMinted - collection.TotalMinted
	bo.OverallAmount = realMintedValue
	bo.AvailableAmount = realMintedValue

	//2. 添加交易记录并更新用户余额
	collection.TotalMinted = newTotalMinted
	collection.MintTimes += 1
	collection.LastBlockHeight = eventInfo.BlockHeight
	collection.LastTransactionId = eventInfo.TransactionId

	return svc.ordHistoricBalancesService.HandleMintInscribeEvent(ctx, eventInfo, collection, bo)

}

func (svc OrdinalIndexerTaskServer) DeployBRC20OrdinalCollection(
	ctx context.Context, eventInfo biz.OrdEventBO) error {

	bo := biz.CollectionBO{
		ChainId:         eventInfo.ChainId,
		Protocol:        eventInfo.Protocol,
		Tick:            eventInfo.Tick,
		EventId:         eventInfo.ID,
		TransactionHash: eventInfo.TransactionHash,
		TransactionId:   eventInfo.TransactionId,
		BlockHeight:     eventInfo.BlockHeight,
		DeployTime:      eventInfo.EventTime,
		CreateTime:      time.Now().Unix(),
		UpdateTime:      time.Now().Unix(),
	}
	// check tick
	uniqueLowerTicker := strings.ToLower(eventInfo.Tick)
	collection, err := svc.collectionService.GetCollection(ctx, eventInfo.ChainId, eventInfo.Protocol, uniqueLowerTicker)
	if err != nil {
		return err
	}
	if collection.Tick == uniqueLowerTicker {
		// dup ticker
		fmt.Printf("ProcessBRC20Deploy deploy event %v, but ticker repeat\n", eventInfo)
		svc.ordinalEventService.UpdateEventStatus(ctx, eventInfo.ID, eventInfo.CollectionId, eventInfo.EventStatus, "invalid")
		return nil
	}
	bo.Tick = uniqueLowerTicker
	eventInfo.Tick = uniqueLowerTicker
	eventInfo.CollectionId = collection.ID

	if eventInfo.Max == "" { // without max
		fmt.Printf("ProcessBRC20Deploy event %v, but max missing\n", eventInfo)
		svc.ordinalEventService.UpdateEventStatus(ctx, eventInfo.ID, eventInfo.CollectionId, eventInfo.EventStatus, "invalid")
		return nil
	}

	// dec
	if dec, err := strconv.ParseUint(eventInfo.Decimal, 10, 64); err != nil || dec > 18 {
		// dec invalid
		fmt.Printf("ProcessBRC20Deploy event %v, but dec invalid. dec: %s\n",
			eventInfo,
			eventInfo.Decimal,
		)
		svc.ordinalEventService.UpdateEventStatus(ctx, eventInfo.ID, eventInfo.CollectionId, eventInfo.EventStatus, "invalid")
		return nil
	} else {
		bo.Decimal = int64(dec)
	}

	// max
	if maxVal, precision, parseErr := decimal.NewDecimalFromString(eventInfo.Max); parseErr != nil {
		// max invalid
		fmt.Printf("ProcessBRC20Deploy event %v, but max invalid. max: '%s', error: %s\n",
			eventInfo,
			eventInfo.Max,
			parseErr.Error(),
		)
		svc.ordinalEventService.UpdateEventStatus(ctx, eventInfo.ID, eventInfo.CollectionId, eventInfo.EventStatus, "invalid")
		return nil
	} else {
		if maxVal.Sign() <= 0 || maxVal.IsOverflowUint64() || precision > int(bo.Decimal) {
			fmt.Printf("ProcessBRC20Deploy event %v, but max precision out of precision. precision: %d, decimal: %d\n",
				eventInfo,
				precision,
				int(bo.Decimal))
			svc.ordinalEventService.UpdateEventStatus(ctx, eventInfo.ID, eventInfo.CollectionId, eventInfo.EventStatus, "invalid")
			return nil
		}
		bo.MaxAmount = maxVal.Float64()
	}

	// lim
	if lim, precision, parseErr := decimal.NewDecimalFromString(eventInfo.Limit); parseErr != nil {
		// limit invalid
		fmt.Printf("ProcessBRC20Deploy event %v, but limit invalid. limit: '%s', error: %s\n",
			eventInfo,
			eventInfo.Limit,
			parseErr.Error(),
		)
		svc.ordinalEventService.UpdateEventStatus(ctx, eventInfo.ID, eventInfo.CollectionId, eventInfo.EventStatus, "invalid")
		return nil
	} else {
		if lim.Sign() <= 0 || lim.IsOverflowUint64() || precision > int(bo.Decimal) {
			fmt.Printf("ProcessBRC20Deploy event %v, but lim precision out of precision. precision: %d, decimal: %d\n",
				eventInfo,
				precision,
				int(bo.Decimal))
			svc.ordinalEventService.UpdateEventStatus(ctx, eventInfo.ID, eventInfo.CollectionId, eventInfo.EventStatus, "invalid")
			return nil
		}
		bo.MintLimitAmount = lim.Float64()
	}

	return svc.ordHistoricBalancesService.HandleDeployInscribeEvent(ctx, eventInfo, bo)

}

func (svc OrdinalIndexerTaskServer) TransferBRC20OrdinalCollection(
	ctx context.Context, eventInfo biz.OrdEventBO) error {

	bo := biz.OrdHistoricBalancesBO{
		ChainId:         eventInfo.ChainId,
		EventId:         eventInfo.ID,
		TransactionHash: eventInfo.TransactionHash,
		TransactionId:   eventInfo.TransactionId,
		BlockHeight:     eventInfo.BlockHeight,
		CreateTime:      time.Now().Unix(),
		UpdateTime:      time.Now().Unix(),
	}

	var transferAmountValue float64
	// check tick
	uniqueLowerTicker := strings.ToLower(eventInfo.Tick)
	collection, err := svc.collectionService.GetCollection(ctx, eventInfo.ChainId, eventInfo.Protocol, uniqueLowerTicker)
	if err != nil {
		return err
	}
	if collection.Tick != uniqueLowerTicker {
		//ticker not exist
		fmt.Printf("TransferBRC20OrdinalCollection transfer event %v, but ticker not exist\n", eventInfo)
		svc.ordinalEventService.UpdateEventStatus(ctx, eventInfo.ID, eventInfo.CollectionId, eventInfo.EventStatus, "invalid")
		return nil
	}
	bo.Tick = uniqueLowerTicker
	bo.CollectionId = collection.ID
	eventInfo.CollectionId = collection.ID
	// check transfer amount
	if amount, precision, parseErr := decimal.NewDecimalFromString(eventInfo.Amount); parseErr != nil {
		// max invalid
		fmt.Printf("TransferBRC20OrdinalCollection event %v, but amount parse err. amount: '%s', error: %s\n",
			eventInfo,
			eventInfo.Amount,
			parseErr.Error(),
		)
		svc.ordinalEventService.UpdateEventStatus(ctx, eventInfo.ID, eventInfo.CollectionId, eventInfo.EventStatus, "invalid")
		return nil
	} else {
		if amount.Sign() <= 0 || precision > int(collection.Decimal) {
			fmt.Printf("TransferBRC20OrdinalCollection event %v, but amount invalid.\n",
				eventInfo)
			svc.ordinalEventService.UpdateEventStatus(ctx, eventInfo.ID, eventInfo.CollectionId, eventInfo.EventStatus, "invalid")
			return nil
		}
		transferAmountValue = amount.Float64()
	}

	//2. 添加交易记录并更新用户余额
	if eventInfo.EventType == biz.TransferInscribe {
		//check if available balance is enough
		OrdHistoricBalancesBO, err := svc.ordHistoricBalancesService.GetLastBalance(
			ctx, eventInfo.ChainId, eventInfo.SourceAddress, bo.CollectionId)
		if err != nil {
			return err
		}
		if transferAmountValue > OrdHistoricBalancesBO.AvailableBalance {
			fmt.Printf("TransferBRC20OrdinalCollection event %v, but transfer inscribe amount %s out of available balance %f in address: %s.\n",
				eventInfo, eventInfo.Amount, OrdHistoricBalancesBO.AvailableBalance, OrdHistoricBalancesBO.WalletAddress)
			svc.ordinalEventService.UpdateEventStatus(ctx, eventInfo.ID, eventInfo.CollectionId, eventInfo.EventStatus, "invalid")
			return nil
		}
		//transfer-inscribe, available_balance -= amount
		bo.AvailableAmount = transferAmountValue
		bo.OverallAmount = 0
		bo.ChangeType = 0 //支出
		bo.WalletAddress = eventInfo.SourceAddress
		bo.Pkscript = eventInfo.SourcePkscript
		err = svc.ordHistoricBalancesService.HandleTransferInscribeEvent(ctx, eventInfo, bo)
		return err
	} else if eventInfo.EventType == biz.TransferTransfer {

		sentAsFeeTransfer := eventInfo.TargetAddress == "" && eventInfo.TargetPkscript == ""

		if sentAsFeeTransfer {
			// sentAsFeeTransfer, sourceAddress available_balance += amount
			bo.OverallAmount = 0
			bo.AvailableAmount = transferAmountValue
			bo.ChangeType = 1 //支出
			bo.WalletAddress = eventInfo.SourceAddress
			bo.Pkscript = eventInfo.SourcePkscript
			err = svc.ordHistoricBalancesService.HandleSentAsFeeTransferEvent(ctx, eventInfo, bo)
			return err

		} else {
			//normal transfer-transfer event, sourceAddress overall_balance -= amount;
			//targetAddress overall_balance += amount, available_balance += amount
			bo.OverallAmount = transferAmountValue
			bo.AvailableAmount = 0
			bo.ChangeType = 0 //支出
			bo.WalletAddress = eventInfo.SourceAddress
			bo.Pkscript = eventInfo.SourcePkscript
			err = svc.ordHistoricBalancesService.HandleTransferTransferEvent(ctx, eventInfo, bo)
			return err

		}

	}

	return nil
}
