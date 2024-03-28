package service

import (
	"context"
	"fmt"
	"github.com/6block/fox_ordinal/pkg/utils"
	"time"
)

func (svc OrdinalIndexerTaskServer) SyncOrdinalEventJob() {
	currentTimeStr := utils.TimestampUnixToFormat(time.Now().Unix(), nil)
	fmt.Printf("【%s】SyncOrdinalEventJob start-----------------------\n", currentTimeStr)
	err := svc.ordinalEventService.SyncEvent(context.Background())
	currentTimeStr = utils.TimestampUnixToFormat(time.Now().Unix(), nil)

	if err != nil {
		fmt.Printf("【%s】SyncOrdinalEventJob failed, error: %s\n", currentTimeStr, err.Error())
	} else {
		fmt.Printf("【%s】SyncOrdinalEventJob success! -----------------------\n", currentTimeStr)
	}
}
